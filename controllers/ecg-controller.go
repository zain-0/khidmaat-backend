package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zain-0/khidmaat-backend/config"
	"github.com/zain-0/khidmaat-backend/models"
	"github.com/zain-0/khidmaat-backend/utils"
	"go.mongodb.org/mongo-driver/bson"
)

// Endpoint to handle the request and send data to Lambda API
func DenoiseData(c *gin.Context) {
	// Declare a variable to hold the incoming request body
	var requestBody utils.SignalRequest

	// Parse the JSON from the incoming request body
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		// If parsing fails, respond with a 400 Bad Request error
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Ensure that signal data is provided
	if len(requestBody.Signal) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Signal data is required"})
		return
	}

	// Call the function to send the data to the Lambda API
	response, err := utils.SendSignalToDenoiseAPI(requestBody.Signal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send signal to API"})
		return
	}

	// Send the response back to the user
	c.JSON(http.StatusOK, gin.H{
		"message": "Signal processed successfully",
		"data":    response,
	})
}

func DetectRPeaksHandler(c *gin.Context) {
	var req utils.SignalRequest

	// Bind the incoming JSON request to the SignalRequest struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Ensure that the signal data is present
	if len(req.Signal) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Signal data is required"})
		return
	}

	// Step 1: Detect R-peaks
	rPeaks := utils.DetectRPeaks(req.Signal)

	// Step 2: Return R-locations and count of detected R-peaks
	c.JSON(http.StatusOK, gin.H{
		"message":    "R-peaks detected successfully",
		"Rlocation":  rPeaks,
		"peak_count": len(rPeaks),
	})
}

func SegmentHeartbeatsHandler(c *gin.Context) {
	var req utils.SignalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Step 1: Denoise the signal
	denoiseResp, err := utils.SendSignalToDenoiseAPI(req.Signal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to denoise signal"})
		return
	}

	// Convert the denoised response to []float64
	rawData, ok := denoiseResp["denoised"]
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No 'denoised' data found in response"})
		return
	}

	interfaceSlice, ok := rawData.([]interface{})
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "'denoised' data is not a list"})
		return
	}

	denoised := make([]float64, len(interfaceSlice))
	for i, v := range interfaceSlice {
		if f, ok := v.(float64); ok {
			denoised[i] = f
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid data in denoised array"})
			return
		}
	}

	// Step 2: Detect R-peaks
	rPeaks := utils.DetectRPeaks(denoised)

	// Step 3: Extract segments from R-99 to R+200
	var segments [][]float64
	for _, r := range rPeaks {
		start := r - 99
		end := r + 200

		if start >= 0 && end < len(denoised) {
			segment := denoised[start:end]
			segments = append(segments, segment)
		}
	}

	// Return the heartbeat segments
	c.JSON(http.StatusOK, gin.H{"segments": segments})
}

func ClassifyHeartbeatsHandler(c *gin.Context) {
	var req utils.SignalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the denoising API
	denoiseResp, err := utils.SendSignalToDenoiseAPI(req.Signal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to denoise signal"})
		return
	}

	// Extract []float64 from response
	rawData, ok := denoiseResp["denoised"]
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No 'denoised' data found in response"})
		return
	}

	interfaceSlice, ok := rawData.([]interface{})
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "'denoised' data is not a list"})
		return
	}

	denoised := make([]float64, len(interfaceSlice))
	for i, v := range interfaceSlice {
		if f, ok := v.(float64); ok {
			denoised[i] = f
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid data in denoised array"})
			return
		}
	}

	// Detect R-peaks
	rPeaks := utils.DetectRPeaks(denoised)
	var results []map[string]interface{}

	for _, r := range rPeaks {
		start := r - 99
		end := r + 200
		if start >= 0 && end < len(denoised) {
			segment := denoised[start:end]

			payload, _ := json.Marshal(utils.PredictionRequest{Segment: segment})

			resp, err := http.Post(
				"https://niml49ck97.execute-api.us-east-1.amazonaws.com/khidmaat-function/predict",
				"application/json",
				bytes.NewBuffer(payload),
			)
			if err != nil {
				results = append(results, map[string]interface{}{
					"error": err.Error(),
				})
				continue
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				results = append(results, map[string]interface{}{
					"error": "Failed to read response body",
				})
				continue
			}

			fmt.Println("API Response Body:", string(body))

			var predictionRes utils.PredictionResponse
			if err := json.Unmarshal(body, &predictionRes); err != nil {
				results = append(results, map[string]interface{}{
					"error": "Invalid prediction response",
				})
				continue
			}

			pred := predictionRes.Prediction[0]
			maxIdx := 0
			for i := 1; i < len(pred); i++ {
				if pred[i] > pred[maxIdx] {
					maxIdx = i
				}
			}

			results = append(results, map[string]interface{}{
				"prediction":  pred,
				"class_label": utils.ClassLabels[maxIdx],
			})

			time.Sleep(200 * time.Millisecond)
		}
	}

	c.JSON(http.StatusOK, gin.H{"classified_segments": results})
}

func FullECGProcessingHandler(c *gin.Context) {
	var requestBody utils.SignalRequest

	// Step 1: Receive raw signal
	if err := c.ShouldBindJSON(&requestBody); err != nil || len(requestBody.Signal) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or empty signal data"})
		return
	}

	// Step 2: Send to denoising Lambda
	denoiseResp, err := utils.SendSignalToDenoiseAPI(requestBody.Signal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to denoise signal"})
		return
	}

	rawDenoised, ok := denoiseResp["denoised_signal"].([]interface{})
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid denoised signal format"})
		return
	}

	// Step 3: Convert denoised signal back to []float64
	denoised := make([]float64, len(rawDenoised))
	for i, v := range rawDenoised {
		f, ok := v.(float64)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid float conversion"})
			return
		}
		denoised[i] = f
	}

	// Step 4: Detect R-peaks
	rPeaks := utils.DetectRPeaks(denoised)

	// Step 5: Segment and classify
	var results []map[string]interface{}
	for _, r := range rPeaks {
		start := r - 99
		end := r + 200
		if start >= 0 && end < len(denoised) {
			segment := denoised[start:end]
			payload, _ := json.Marshal(utils.PredictionRequest{Segment: segment})

			resp, err := http.Post(
				"https://niml49ck97.execute-api.us-east-1.amazonaws.com/khidmaat-function/predict",
				"application/json",
				bytes.NewBuffer(payload),
			)
			if err != nil {
				results = append(results, map[string]interface{}{"error": err.Error()})
				continue
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				results = append(results, map[string]interface{}{"error": "Failed to read response body"})
				continue
			}

			var predictionRes utils.PredictionResponse
			if err := json.Unmarshal(body, &predictionRes); err != nil {
				results = append(results, map[string]interface{}{"error": "Invalid prediction response"})
				continue
			}

			// ✅ Safe check added here
			if len(predictionRes.Prediction) == 0 {
				results = append(results, map[string]interface{}{"error": "Empty prediction array"})
				continue
			}

			pred := predictionRes.Prediction[0]
			if len(pred) == 0 {
				results = append(results, map[string]interface{}{"error": "Empty probability values in prediction"})
				continue
			}

			maxIdx := 0
			for i := 1; i < len(pred); i++ {
				if pred[i] > pred[maxIdx] {
					maxIdx = i
				}
			}

			results = append(results, map[string]interface{}{
				"r_peak_index": r,
				"class_label":  utils.ClassLabels[maxIdx],
				"confidence":   pred[maxIdx],
				"raw_probs":    pred,
			})

			time.Sleep(200 * time.Millisecond)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":             "ECG signal processed and classified",
		"r_peak_count":        len(rPeaks),
		"classified_segments": results,
	})
}

func AlertECGHandler(c *gin.Context) {
	var requestBody struct {
		UserID string    `json:"user_id"`
		Signal []float64 `json:"signal"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil || len(requestBody.Signal) == 0 || requestBody.UserID == "" {
		log.Println("❌ Invalid request body:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request, missing user_id or signal"})
		return
	}
	log.Printf("📥 Received ECG signal for user_id: %s, length: %d\n", requestBody.UserID, len(requestBody.Signal))

	// Fetch hospital_id from users collection
	var user struct {
		UserID     string `bson:"user_id"`
		HospitalID string `bson:"hospital_id"`
	}
	if err := config.UsersCollection.FindOne(context.TODO(), bson.M{"user_id": requestBody.UserID}).Decode(&user); err != nil {
		log.Printf("❌ User not found for user_id: %s. Error: %v\n", requestBody.UserID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	log.Printf("🏥 Retrieved hospital_id: %s for user_id: %s\n", user.HospitalID, user.UserID)

	// Denoising
	denoiseResp, err := utils.SendSignalToDenoiseAPI(requestBody.Signal)
	if err != nil {
		log.Println("❌ Failed to denoise signal:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to denoise signal"})
		return
	}

	rawDenoised, ok := denoiseResp["denoised_signal"].([]interface{})
	if !ok {
		log.Println("❌ Denoised signal format invalid")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid denoised signal format"})
		return
	}

	denoised := make([]float64, len(rawDenoised))
	for i, v := range rawDenoised {
		f, ok := v.(float64)
		if !ok {
			log.Printf("❌ Failed to convert index %d to float\n", i)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid float conversion"})
			return
		}
		denoised[i] = f
	}
	log.Println("✅ Denoising complete")

	rPeaks := utils.DetectRPeaks(denoised)
	log.Printf("📈 Detected %d R-peaks\n", len(rPeaks))

	var results []map[string]interface{}
	classCount := make(map[string]int)
	var heartBeats []models.HeartBeat

	for i, r := range rPeaks {
		start := r - 99
		end := r + 200
		if start >= 0 && end < len(denoised) {
			segment := denoised[start:end]
			payload, _ := json.Marshal(utils.PredictionRequest{Segment: segment})

			resp, err := http.Post(
				"https://niml49ck97.execute-api.us-east-1.amazonaws.com/khidmaat-function/predict",
				"application/json",
				bytes.NewBuffer(payload),
			)
			if err != nil {
				log.Printf("❌ Request %d failed to predict: %v\n", i, err)
				continue
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Printf("❌ Failed to read prediction response at index %d: %v\n", i, err)
				continue
			}

			var predictionRes utils.PredictionResponse
			if err := json.Unmarshal(body, &predictionRes); err != nil {
				log.Printf("❌ Failed to unmarshal prediction response at index %d: %v\n", i, err)
				continue
			}

			pred := predictionRes.Prediction[0]
			maxIdx := 0
			for i := 1; i < len(pred); i++ {
				if pred[i] > pred[maxIdx] {
					maxIdx = i
				}
			}

			label := utils.ClassLabels[maxIdx]
			classCount[label]++
			heartBeats = append(heartBeats, models.HeartBeat{
				Label:      label,
				Confidence: pred[maxIdx],
			})

			results = append(results, map[string]interface{}{
				"r_peak_index": r,
				"class_label":  label,
				"confidence":   pred[maxIdx],
				"raw_probs":    pred,
			})

			time.Sleep(200 * time.Millisecond)
		}
	}

	alert := false
	if classCount["Premature ventr."]+classCount["Supraventricular prem"] >= 2 {
		alert = true
	}
	if classCount["Normal"] < (len(rPeaks) - classCount["Normal"]) {
		alert = true
	}
	log.Printf("🚨 Alert determination complete. Alert = %v\n", alert)

	// Store Medical Record
	medicalRecord := models.MedicalRecord{
		UserID:     requestBody.UserID,
		MedicalID:  uuid.New().String(),
		HeartBeats: heartBeats,
		Conclusion: "Alert Needed: " + strconv.FormatBool(alert),
		SentAt:     time.Now().Format(time.RFC3339),
	}

	// Only set HospitalTreated if alert is true
	if alert {
		medicalRecord.HospitalTreated = &user.HospitalID
	}

	insertResult, err := config.MedicalRecordsCollection.InsertOne(context.TODO(), medicalRecord)
	if err != nil {
		log.Printf("❌ Failed to store medical record: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save medical record"})
		return
	}
	log.Printf("✅ Medical record stored. ID: %v\n", insertResult.InsertedID)

	c.JSON(http.StatusOK, gin.H{
		"message":             "Processed and stored medical record",
		"r_peak_count":        len(rPeaks),
		"classified_segments": results,
		"alert_required":      alert,
	})
}
