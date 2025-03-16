package auth

import (
	"drones-be/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authSrv  *services.AuthServices
	tokenSrv *services.TokenServices
}

func NewAuthHandler(authSrv *services.AuthServices, tokenSrv *services.TokenServices) *AuthHandler {
	return &AuthHandler{
		authSrv:  authSrv,
		tokenSrv: tokenSrv,
	}
}

type RegisterUser struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) RegisterUser(c *gin.Context) {

	var req RegisterUser
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authSrv.RegisterUser(req.Username, req.Password, "user")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"user": user})

}

type LoginUser struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) LoginUser(c *gin.Context) {

	var req LoginUser
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authSrv.LoginUser(req.Username, req.Password)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userID := user.ID.String()

	token, err := h.tokenSrv.GenerateToken(userID, user.Role)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("token", token, 60*60*24, "/", "localhost", false, true)

	c.JSON(200, gin.H{"user_id": userID, "username": user.Username, "role": user.Role})

}
