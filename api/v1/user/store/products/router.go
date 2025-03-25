package storeproducts

import (
	"drones-be/internal/middleware"
	"drones-be/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, h *ProductHandler, token *services.TokenServices) {

	storeProductsGroup := r.Group("/store/products")
	storeProductsGroup.Use(middleware.AuthMiddleware(token))

	{
		storeProductsGroup.POST("", h.GetAvailableProducts)
	}
}
