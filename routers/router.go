package routers

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Setup individual route groups
	SetupUserRoutes(router)
	SetupECGRoutes(router)

	return router
}
