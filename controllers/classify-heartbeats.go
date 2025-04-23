package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil" // Added to read the response body
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var classLabels = []string{
	"Normal",
	"Supraventricular prem",
	"Premature ventr.",
	"Fusion of ve. & no.",
	"Unclassifiable",
}

type PredictionRequest struct {
	Segment []float64 `json:"segment"`
}

type PredictionResponse struct {
	Prediction [][]float64 `json:"prediction"` // Change this to an array of arrays
}

func ClassifyHeartbeatsHandler(c *gin.Context) {
	var req SignalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	denoised := req.Signal
	rPeaks := DetectRPeaks(denoised)

	var results []map[string]interface{}

	for _, r := range rPeaks {
		start := r - 99
		end := r + 200
		if start >= 0 && end < len(denoised) {
			segment := denoised[start:end]

			// Prepare payload
			payload, _ := json.Marshal(PredictionRequest{Segment: segment})

			// Send request to model API
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

			// Read the response body for debugging
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				results = append(results, map[string]interface{}{
					"error": "Failed to read response body",
				})
				continue
			}

			// Log the response body to debug the API response
			// You can remove this line after debugging
			fmt.Println("API Response Body:", string(body))

			var predictionRes PredictionResponse
			if err := json.Unmarshal(body, &predictionRes); err != nil {
				results = append(results, map[string]interface{}{
					"error": "Invalid prediction response",
				})
				continue
			}

			// The prediction is now the first element in the "prediction" array
			pred := predictionRes.Prediction[0] // Get the first array of predictions
			maxIdx := 0
			for i := 1; i < len(pred); i++ {
				if pred[i] > pred[maxIdx] {
					maxIdx = i
				}
			}

			results = append(results, map[string]interface{}{
				"prediction":  pred,
				"class_label": classLabels[maxIdx],
			})

			time.Sleep(200 * time.Millisecond)
		}
	}

	c.JSON(http.StatusOK, gin.H{"classified_segments": results})
}
