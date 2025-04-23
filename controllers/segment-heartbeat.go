package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SegmentHeartbeatsHandler(c *gin.Context) {
	var req SignalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Step 1: Denoise the signal
	denoised := req.Signal

	// Step 2: Detect R-peaks
	rPeaks := DetectRPeaks(denoised)

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
