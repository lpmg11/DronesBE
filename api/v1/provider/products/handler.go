package providerproducts

import (
	"drones-be/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductHandler struct {
	db *gorm.DB
}

func NewProductHandler(db *gorm.DB) *ProductHandler {
	return &ProductHandler{db: db}
}

type NewProduct struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req NewProduct

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Solicitud inválida", "detalles": err.Error()})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		c.JSON(401, gin.H{"error": "Usuario no encontrado"})
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(400, gin.H{"error": "ID de usuario inválido", "detalles": err.Error()})
		return
	}

	provider := models.Provider{}
	if err := h.db.Where("user_id = ?", userUUID).First(&provider).Error; err != nil {
		c.JSON(500, gin.H{"error": "Error al buscar proveedor", "detalles": err.Error()})
		return
	}

	product := models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		ProviderID:  provider.ID,
	}

	if err := h.db.Create(&product).Error; err != nil {
		c.JSON(500, gin.H{"error": "Error al crear producto", "detalles": err.Error()})
		return
	}

	c.JSON(201, gin.H{"product": product})

}

func (h *ProductHandler) GetProducts(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		c.JSON(401, gin.H{"error": "Usuario no encontrado"})
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(400, gin.H{"error": "ID de usuario inválido", "detalles": err.Error()})
		return
	}

	provider := models.Provider{}
	if err := h.db.Where("user_id = ?", userUUID).First(&provider).Error; err != nil {
		c.JSON(500, gin.H{"error": "Error al buscar proveedor", "detalles": err.Error()})
		return
	}

	products := []models.Product{}
	if err := h.db.Where("provider_id = ?", provider.ID).Find(&products).Error; err != nil {
		c.JSON(500, gin.H{"error": "Error al buscar productos", "detalles": err.Error()})
		return
	}

	c.JSON(200, gin.H{"products": products})
}