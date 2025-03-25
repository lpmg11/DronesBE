package admindrone

import (
	"drones-be/internal/middleware"
	"drones-be/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, h *DroneHandler, token *services.TokenServices) {

	droneGroup := r.Group("/drone")
	droneGroup.Use(middleware.AuthMiddleware(token))

	{
		droneGroup.POST("/model", h.CreateDroneModel)
		droneGroup.GET("/model", h.GetDroneModel)
		droneGroup.POST("", h.CreateDrone)
		droneGroup.GET("", h.GetDrones)
	}
}
