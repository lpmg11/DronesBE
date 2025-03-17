package adminwarehouse

import (
	"drones-be/internal/models"
	"drones-be/internal/storage"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WarehouseHandler struct {
	db *gorm.DB
}

func NewWarehouseHandler(db *storage.PostgresClient) *WarehouseHandler {
	return &WarehouseHandler{db: db.DB}
}

type NewWareHouse struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (h *WarehouseHandler) CreateWarehouse(c *gin.Context) {
	var req NewWareHouse

	role, exist := c.Get("role")

	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no encontrado"})
		return
	}

	if role != "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No tienes permisos para realizar esta acción"})
		return
	}

	if err := c.ShouldBindJSON(&req); err !=
		nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos invalidos", "detalles": err.Error()})
		return
	}

	var warehouse models.Warehouse

	warehouse.Name = req.Name
	warehouse.Latitude = req.Latitude
	warehouse.Longitude = req.Longitude

	if err := h.db.Create(&warehouse).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear almacén", "detalles": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"warehouse": warehouse})

}

func (h *WarehouseHandler) GetWarehouses(c *gin.Context) {
	var warehouses []models.Warehouse

	if err := h.db.Find(&warehouses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener almacenes", "detalles": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"warehouses": warehouses})
}

type WarehouseWithDistance struct {
	models.Warehouse
	Distance float64 `json:"distance"`
}

func (h *WarehouseHandler) GetWarehousesByProximity(c *gin.Context) {
	var req struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Radius    float64 `json:"radius"` // en kilómetros
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos invalidos", "detalles": err.Error()})
		return
	}

	var warehouses []WarehouseWithDistance

	query := `
		SELECT 
			id, name, latitude, longitude,
			earth_distance(ll_to_earth(?, ?), ll_to_earth(latitude, longitude)) AS distance
		FROM warehouses
		WHERE earth_box(ll_to_earth(?, ?), ? * 1000) @> ll_to_earth(latitude, longitude)
		ORDER BY distance
	`

	if err := h.db.Raw(query,
		req.Latitude, req.Longitude,
		req.Latitude, req.Longitude, req.Radius).Scan(&warehouses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener almacenes", "detalles": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"warehouses": warehouses})
}
