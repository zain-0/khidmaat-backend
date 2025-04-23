package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/zain-0/khidmaat-backend/controllers"
)

func SetupECGRoutes(router *gin.Engine) {
	router.POST("/denoise-signal", controllers.DenoiseData)
	router.POST("/detect-rpeaks", controllers.DetectRPeaksHandler)
	router.POST("/segment-heartbeats", controllers.SegmentHeartbeatsHandler)
	router.POST("/classify-heartbeats", controllers.ClassifyHeartbeatsHandler)
	router.POST("/process-ecg", controllers.FullECGProcessingHandler)
}
