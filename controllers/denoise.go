package controllers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Define the structure of the request body
type SignalRequest struct {
	Signal []float64 `json:"signal"`
}

// Function to send the signal data to the Lambda API
func sendSignalToDenoiseAPI(signalData []float64) (map[string]interface{}, error) {
	url := "https://9zpm7oy3s2.execute-api.us-east-1.amazonaws.com/khidmaat-denoise"
	// Create the request body
	requestBody := SignalRequest{Signal: signalData}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		return nil, err
	}

	// Send POST request to the API
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error sending POST request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Check if response status is OK
	if resp.StatusCode != http.StatusOK {
		log.Printf("Error: received non-200 status code %v", resp.StatusCode)
		return nil, err
	}

	// Parse the response body (assuming the response is JSON)
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Printf("Error decoding response: %v", err)
		return nil, err
	}

	return response, nil
}

// Endpoint to handle the request and send data to Lambda API
func DenoiseData(c *gin.Context) {
	// Declare a variable to hold the incoming request body
	var requestBody SignalRequest

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
	response, err := sendSignalToDenoiseAPI(requestBody.Signal)
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
