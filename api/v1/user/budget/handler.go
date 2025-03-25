package userbudget

import (
	"drones-be/internal/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BudgetHandler struct {
	db *gorm.DB
}

func NewBudgetHandler(db *gorm.DB) *BudgetHandler {
	return &BudgetHandler{db: db}
}

func (h *BudgetHandler) GetBudget(c *gin.Context) {
	// Obtener el userID del contexto
	user, exist := c.Get("userID")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no encontrado"})
		return
	}

	userID, err := uuid.Parse(user.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuario inválido", "detalles": err.Error()})
		return
	}

	// Buscar el cliente asociado al usuario
	var client models.Client
	err = h.db.Where("user_id = ?", userID).First(&client).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Si no existe, se crea el cliente usando latitud y longitud 0 y el username del usuario
			var userModel models.User
			if err := h.db.First(&userModel, "id = ?", userID).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener datos del usuario", "detalles": err.Error()})
				return
			}

			client = models.Client{
				Name:      userModel.Username,
				UserID:    userID,
				Latitude:  0,
				Longitude: 0,
			}
			if err := h.db.Create(&client).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear cliente", "detalles": err.Error()})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar cliente", "detalles": err.Error()})
			return
		}
	}

	// Buscar el presupuesto asociado al cliente
	var budget models.Budget
	err = h.db.Preload("Transactions", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC")
	}).Where("client_id = ?", client.ID).First(&budget).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No existe presupuesto, por lo que se crea uno nuevo
			budget = models.Budget{
				ClientID: client.ID,
				Balance:  0,
			}
			if err := h.db.Create(&budget).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear presupuesto", "detalles": err.Error()})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar presupuesto", "detalles": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"budget": budget})
}

func (h *BudgetHandler) CreateBudget(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no encontrado"})
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuario inválido", "detalles": err.Error()})
		return
	}

	var budget models.Budget
	budget.ClientID = userUUID
	budget.Balance = 0

	if err := h.db.Create(&budget).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear presupuesto", "detalles": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"budget": budget})
}

// FoundRequestType es la estructura para recibir la solicitud de fondos.
type FoundRequestType struct {
	Amount           float64 `json:"amount"`
	Description      string  `json:"description"`
	ConfirmationCode string  `json:"confirmation_code"`
}

func (h *BudgetHandler) FoundRequest(c *gin.Context) {
	var req FoundRequestType

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Solicitud inválida", "detalles": err.Error()})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no encontrado"})
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuario inválido", "detalles": err.Error()})
		return
	}

	var client models.Client
	err = h.db.Where("user_id = ?", userUUID).First(&client).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar cliente", "detalles": err.Error()})
		return
	}

	var budget models.Budget
	err = h.db.Where("client_id = ?", client.ID).First(&budget).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar presupuesto", "detalles": err.Error()})
		return
	}

	transaction := models.BudgetTransaction{
		Amount:           req.Amount,
		Description:      req.Description,
		ConfirmationCode: req.ConfirmationCode,
		BudgetID:         budget.ID,
		Status:           "Pendiente",
	}

	if err := h.db.Create(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear transacción", "detalles": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"transaction": transaction})
}
