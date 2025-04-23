package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/zain-0/khidmaat-backend/controllers"
)

func SetupHospitalRoutes(router *gin.Engine) {
	hospitalGroup := router.Group("/api/hospitals")
	{
		hospitalGroup.POST("/", controllers.CreateHospital)
		hospitalGroup.GET("/", controllers.GetAllHospitals)
		hospitalGroup.GET("/:id", controllers.GetHospitalByID)
		hospitalGroup.PUT("/:id", controllers.UpdateHospital)
		hospitalGroup.DELETE("/:id", controllers.DeleteHospital)
	}
}
