package controllers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// --- Final Combined Flow Handler ---

func FullECGProcessingHandler(c *gin.Context) {
	var requestBody SignalRequest

	// Step 1: Receive raw signal
	if err := c.ShouldBindJSON(&requestBody); err != nil || len(requestBody.Signal) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or empty signal data"})
		return
	}

	// Step 2: Send to denoising Lambda
	denoiseResp, err := sendSignalToDenoiseAPI(requestBody.Signal)
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
	rPeaks := DetectRPeaks(denoised)

	// Step 5: Segment and classify
	var results []map[string]interface{}
	for _, r := range rPeaks {
		start := r - 99
		end := r + 200
		if start >= 0 && end < len(denoised) {
			segment := denoised[start:end]
			payload, _ := json.Marshal(PredictionRequest{Segment: segment})

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

			var predictionRes PredictionResponse
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
				"class_label":  classLabels[maxIdx],
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
