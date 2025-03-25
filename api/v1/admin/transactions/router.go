package admintransactions

import (
	"drones-be/internal/middleware"
	"drones-be/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, h *TransactionHandler, token *services.TokenServices) {
	adminTransactionsGroup := r.Group("/admin/transactions")
	adminTransactionsGroup.Use(middleware.AuthMiddleware(token))

	{
		adminTransactionsGroup.GET("", h.GetTransactions)
		adminTransactionsGroup.PUT("", h.AproveTransaction)
	}
}
