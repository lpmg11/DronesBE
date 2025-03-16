package auth

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, h *AuthHandler) {
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", h.RegisterUser)
		authGroup.POST("/login", h.LoginUser)
		authGroup.POST("/logout", h.LogoutUser)
	}
}
