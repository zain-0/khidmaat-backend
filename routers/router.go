package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/zain-0/khidmaat-backend/controllers"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// User routes
	router.POST("/signup", controllers.SignUp)
	router.POST("/login", controllers.Login)
	router.GET("/user/:id", controllers.GetUserWithDetails)
	router.GET("/users", controllers.GetUsersByQuery)

	// Add the new route for processing signal data
	router.POST("/denoise-signal", controllers.ProcessSignalData)
	router.POST("/detect-rpeaks", controllers.DetectRPeaksHandler)

	// Add other routes here

	return router
}
