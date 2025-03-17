package admindrone

import (
	"drones-be/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DroneHandler struct {
	db *gorm.DB
}

func NewDroneHandler(db *gorm.DB) *DroneHandler {
	return &DroneHandler{db: db}
}

type NewDroneModel struct {
	Name     string  `json:"name"`
	ChargeKM float64 `json:"charge_km"`
	Speed    float64 `json:"speed"`
}

func (h *DroneHandler) CreateDroneModel(c *gin.Context) {

	var req NewDroneModel

	verifyIsAdmin(c)

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Solicitud inválida", "detalles": err.Error()})
		return
	}

	droneModel := &models.DroneModel{
		Name:     req.Name,
		ChargeKM: req.ChargeKM,
		Speed:    req.Speed,
	}

	if err := h.db.Create(droneModel).Error; err != nil {
		c.JSON(500, gin.H{"error": "Error al crear modelo de dron", "detalles": err.Error()})
		return
	}

	c.JSON(200, gin.H{"drone_model": droneModel})

}

func (h *DroneHandler) GetDroneModel(c *gin.Context) {
	var droneModels []models.DroneModel

	verifyIsAdmin(c)

	if err := h.db.Find(&droneModels).Error; err != nil {
		c.JSON(500, gin.H{"error": "Error al obtener modelos de drones", "detalles": err.Error()})
		return
	}

	c.JSON(200, gin.H{"drone_models": droneModels})
}

type NewDrone struct {
	Name        string `json:"name"`
	WarehouseID string `json:"warehouse_id"`
	ModelID     string `json:"model_id"`
}

func (h *DroneHandler) CreateDrone(c *gin.Context) {
	var req NewDrone

	verifyIsAdmin(c)

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Solicitud inválida", "detalles": err.Error()})
		return
	}

	warehouseID, err := uuid.Parse(req.WarehouseID)
	if err != nil {
		c.JSON(400, gin.H{"error": "ID de almacén inválido", "detalles": err.Error()})
		return
	}

	modelID, err := uuid.Parse(req.ModelID)
	if err != nil {
		c.JSON(400, gin.H{"error": "ID de modelo inválido", "detalles": err.Error()})
		return
	}

	drone := &models.Drone{
		Name:        req.Name,
		WarehouseID: warehouseID,
		ModelId:     modelID,
	}

	if err := h.db.Create(drone).Error; err != nil {
		c.JSON(500, gin.H{"error": "Error al crear dron", "detalles": err.Error()})
		return
	}

}

func (h *DroneHandler) GetDrones(c *gin.Context) {
	var drones []models.Drone

	verifyIsAdmin(c)

	if err := h.db.Find(&drones).Error; err != nil {
		c.JSON(500, gin.H{"error": "Error al obtener drones", "detalles": err.Error()})
		return
	}

	c.JSON(200, gin.H{"drones": drones})
}

func verifyIsAdmin(c *gin.Context) {
	role, exist := c.Get("role")

	if !exist {
		c.JSON(401, gin.H{"error": "Usuario no encontrado"})
		return
	}

	if role != "admin" {
		c.JSON(403, gin.H{"error": "No tienes permisos para realizar esta acción"})
	}
}
