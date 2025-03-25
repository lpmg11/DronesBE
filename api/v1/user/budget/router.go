package userbudget

import (
	"drones-be/internal/middleware"
	"drones-be/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, h *BudgetHandler, token *services.TokenServices) {
	userBudgetGroup := r.Group("/user/budget")
	userBudgetGroup.Use(middleware.AuthMiddleware(token))

	{
		userBudgetGroup.POST("", h.CreateBudget)
		userBudgetGroup.GET("", h.GetBudget)
		userBudgetGroup.POST("/request", h.FoundRequest)
	}
}
