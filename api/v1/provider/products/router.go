package providerproducts

import (
	"drones-be/internal/middleware"
	"drones-be/internal/services"

	"github.com/gin-gonic/gin"
)


func RegisterRoutes(r *gin.RouterGroup, h *ProductHandler, token *services.TokenServices) {

	providerProductsGroup := r.Group("/provider/products")
	providerProductsGroup.Use(middleware.AuthMiddleware(token))

	{
		providerProductsGroup.POST("", h.CreateProduct)
		providerProductsGroup.GET("", h.GetProducts)
	}
}