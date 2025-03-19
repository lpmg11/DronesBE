package auth

import (
	"drones-be/internal/config"
	"drones-be/internal/services"
	"log"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authSrv  *services.AuthServices
	tokenSrv *services.TokenServices
	cfg      *config.Config
}

func NewAuthHandler(authSrv *services.AuthServices, tokenSrv *services.TokenServices, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		authSrv:  authSrv,
		tokenSrv: tokenSrv,
		cfg:      cfg,
	}
}

type RegisterUser struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) RegisterUser(c *gin.Context) {

	var req RegisterUser
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Solicitud inválida", "detalles": err.Error()})
		return
	}

	user, err := h.authSrv.RegisterUser(req.Username, req.Password, "admin")
	log.Print(err)
	if err != nil {
		c.JSON(400, gin.H{"error": "Error al registrar usuario", "detalles": err.Error()})
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
		c.JSON(400, gin.H{"error": "Solicitud inválida", "detalles": err.Error()})
		return
	}

	user, err := h.authSrv.LoginUser(req.Username, req.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": "Usuario o contraseña inválidos", "detalles": err.Error()})
		return
	}

	userID := user.ID.String()

	token, err := h.tokenSrv.GenerateToken(userID, user.Role)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error al generar el token", "detalles": err.Error()})
		return
	}

	if h.cfg.Environtment == "production" {
		c.SetCookie("token", token, 60*60*24, "/", h.cfg.Domain, true, true)
	} else {
		c.SetCookie("token", token, 60*60*24, "/", "localhost", false, true)
	}

	c.JSON(200, gin.H{"username": user.Username, "role": user.Role})

}

func (h *AuthHandler) LogoutUser(c *gin.Context) {

	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.JSON(200, gin.H{"message": "Sesión cerrada exitosamente"})

}
