package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/zain-0/khidmaat-backend/controllers"
)

func SetupUserRoutes(router *gin.Engine) {
	userGroup := router.Group("/api/users")
	{
		userGroup.POST("/signup", controllers.SignUp)
		userGroup.POST("/login", controllers.Login)
		userGroup.GET("/:id", controllers.GetUserWithDetails)
		userGroup.GET("/", controllers.GetUsersByQuery)
	}
}
