package utils

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

// Function to send the signal data to the Lambda API
func SendSignalToDenoiseAPI(signalData []float64) (map[string]interface{}, error) {
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
