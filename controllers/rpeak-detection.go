package controllers

import (
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Bandpass filter function (placeholder)
func bandpassFilter(signal []float64, fs float64) []float64 {
	// A basic bandpass filter implementation (for demonstration purposes)
	// Normally, you would use a proper filter design (e.g., Butterworth, Chebyshev).
	// Here we simulate a simple low-pass and high-pass filter.

	nyquist := 0.5 * fs
	_ = 0.5 / nyquist
	_ = 50.0 / nyquist

	// Example simple filter, this is where a real filter implementation should go.
	// This is just a placeholder for an actual bandpass filter.
	filtered := make([]float64, len(signal))
	copy(filtered, signal)
	return filtered // Returning the signal without changes for now.
}

func differentiate(signal []float64) []float64 {
	differentiated := make([]float64, len(signal)-1)
	for i := 1; i < len(signal); i++ {
		differentiated[i-1] = signal[i] - signal[i-1]
	}
	return differentiated
}

func square(signal []float64) []float64 {
	squared := make([]float64, len(signal))
	for i, v := range signal {
		squared[i] = v * v
	}
	return squared
}

func movingWindowIntegration(signal []float64, fs float64) []float64 {
	windowSize := int(0.150 * fs)
	if windowSize < 1 {
		windowSize = 1
	}
	mwi := make([]float64, len(signal))
	for i := range signal {
		start := int(math.Max(float64(i-windowSize/2), 0))
		end := int(math.Min(float64(i+windowSize/2), float64(len(signal)-1)))
		windowSum := 0.0
		for j := start; j <= end; j++ {
			windowSum += signal[j]
		}
		mwi[i] = windowSum / float64(end-start+1)
	}
	return mwi
}

func findPeaks(signal []float64, threshold float64, minDistance int) []int {
	peaks := []int{}
	for i := 1; i < len(signal)-1; i++ {
		if signal[i] > signal[i-1] && signal[i] > signal[i+1] && signal[i] > threshold {
			if len(peaks) == 0 || i-peaks[len(peaks)-1] > minDistance {
				peaks = append(peaks, i)
			}
		}
	}
	return peaks
}

func max(slice []float64) float64 {
	maxValue := slice[0]
	for _, v := range slice[1:] {
		if v > maxValue {
			maxValue = v
		}
	}
	return maxValue
}

// DetectRPeaks detects the R-peaks in the given ECG signal
func DetectRPeaks(signal []float64) []int {
	fs := 360.0
	// Step 1: Bandpass filter the signal (this will be where you apply your real filter)
	filtered := bandpassFilter(signal, fs)

	// Step 2: Differentiate the signal
	diff := differentiate(filtered)

	// Step 3: Square the differentiated signal
	squared := square(diff)

	// Step 4: Moving window integration
	mwi := movingWindowIntegration(squared, fs)

	// Step 5: Apply a threshold to the MWI signal
	threshold := 0.5 * max(mwi)

	// Step 6: Find the peaks in the MWI signal
	minDistance := int(fs / 2) // Minimum distance between peaks (in samples)
	peaks := findPeaks(mwi, threshold, minDistance)

	return peaks
}

// Controller for ECG signal processing
func DetectRPeaksHandler(c *gin.Context) {
	var req SignalRequest

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

	// Step 1: Denoise signal using bandpass filter (the filtered signal is calculated here)
	fs := 360.0
	filteredSignal := bandpassFilter(req.Signal, fs)

	// Step 2: Detect R-peaks
	rPeaks := DetectRPeaks(filteredSignal)

	// Step 3: Return R-locations and count of detected R-peaks
	c.JSON(http.StatusOK, gin.H{
		"message":    "R-peaks detected successfully",
		"Rlocation":  rPeaks,      // R-peak locations (indices of detected peaks)
		"peak_count": len(rPeaks), // Count of detected R-peaks
	})
}
