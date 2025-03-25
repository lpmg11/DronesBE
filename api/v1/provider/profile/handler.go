package providerprofile

import (
	"drones-be/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProfileHandler struct {
	db *gorm.DB
}

func NewProfileHandler(db *gorm.DB) *ProfileHandler {
	return &ProfileHandler{db: db}
}

type NewProfile struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (h *ProfileHandler) CreateProfile(c *gin.Context) {
	var req NewProfile

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Solicitud inválida", "detalles": err.Error()})
		return
	}

	user, exist := c.Get("userID")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no encontrado"})
		return
	}

	userID, err := uuid.Parse(user.(string))
	if err != nil {
		c.JSON(400, gin.H{"error": "ID de usuario inválido", "detalles": err.Error()})
		return
	}

	// Verificar si ya existe un perfil
	var count int64
	if err := h.db.Model(&models.Provider{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		c.JSON(500, gin.H{"error": "Error al verificar perfil", "detalles": err.Error()})
		return
	}

	if count > 0 {
		c.JSON(400, gin.H{"error": "Perfil ya existe", "mensaje": "Ya tienes un perfil de proveedor"})
		return
	}

	// Verificar si hay al menos una bodega a menos de 10 km (10000 metros)
	var nearbyWarehouses int64
	query := `
		SELECT COUNT(*) 
		FROM warehouses
		WHERE earth_box(ll_to_earth(?, ?), 10000) @> ll_to_earth(latitude, longitude)
		  AND earth_distance(ll_to_earth(?, ?), ll_to_earth(latitude, longitude)) <= 10000
	`
	if err := h.db.Raw(query, req.Latitude, req.Longitude, req.Latitude, req.Longitude).Scan(&nearbyWarehouses).Error; err != nil {
		c.JSON(500, gin.H{"error": "Error al verificar proximidad de bodegas", "detalles": err.Error()})
		return
	}

	if nearbyWarehouses == 0 {
		c.JSON(400, gin.H{
			"error":   "No hay bodegas cercanas",
			"mensaje": "Actualmente no hay bodegas a menos de 10 km de tu ubicación. Por favor, espera a que abramos más sucursales cerca de tu oficina.",
		})
		return
	}

	profile := &models.Provider{
		UserID:    userID,
		Name:      req.Name,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	}

	if err := h.db.Create(profile).Error; err != nil {
		c.JSON(500, gin.H{"error": "Error al crear perfil", "detalles": err.Error()})
		return
	}

	c.JSON(200, gin.H{"profile": profile})
}

func (h *ProfileHandler) GetProfile(c *gin.Context) {

	User, exist := c.Get("userID")

	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no encontrado"})
		return
	}

	Role, exist := c.Get("role")

	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Rol no encontrado"})
		return
	}

	if Role != "provider" {
		c.JSON(http.StatusForbidden, gin.H{"error": "No autorizado"})
		return
	}

	UserID := uuid.MustParse(User.(string))

	var profile models.Provider

	if err := h.db.Where("user_id = ?", UserID).First(&profile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Error al obtener perfil", "detalles": err.Error()})
		return
	}

	c.JSON(200, gin.H{"profile": profile})

}
