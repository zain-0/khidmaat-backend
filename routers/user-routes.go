package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/zain-0/khidmaat-backend/controllers"
)

func SetupUserRoutes(router *gin.Engine) {
	router.POST("/signup", controllers.SignUp)
	router.POST("/login", controllers.Login)
	router.GET("/user/:id", controllers.GetUserWithDetails)
	router.GET("/users", controllers.GetUsersByQuery)
}
