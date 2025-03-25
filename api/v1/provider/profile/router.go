package providerprofile

import (
	"drones-be/internal/middleware"
	"drones-be/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, h *ProfileHandler, token *services.TokenServices) {

	providerProfileGroup := r.Group("/provider/profile")
	providerProfileGroup.Use(middleware.AuthMiddleware(token))

	{
		providerProfileGroup.POST("", h.CreateProfile)
		providerProfileGroup.GET("", h.GetProfile)
	}
}
