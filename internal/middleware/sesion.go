package middleware

import (
	"drones-be/internal/services"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(tokenService *services.TokenServices) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err != nil {
			c.JSON(401, gin.H{"error": "no autorizado"})
			c.Abort()
			return
		}

		if token == "" {
			c.JSON(401, gin.H{"error": "no autorizado"})
			c.Abort()
			return
		}

		userID, role, err := tokenService.VerifyToken(token)
		if err != nil {
			c.JSON(401, gin.H{"error": "no autorizado"})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Set("role", role)

		c.Next()
	}
}
