package adminwarehouse

import (
	"drones-be/internal/middleware"
	"drones-be/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, h *WarehouseHandler, token *services.TokenServices) {

	warehouseGroup := r.Group("/warehouse")
	warehouseGroup.Use(middleware.AuthMiddleware(token))

	{
		warehouseGroup.POST("/", h.CreateWarehouse)
		warehouseGroup.GET("/", h.GetWarehouses)
		warehouseGroup.POST("/proximity", h.GetWarehousesByProximity)
	}
}
