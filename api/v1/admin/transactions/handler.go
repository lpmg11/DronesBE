package admintransactions

import (
	"drones-be/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TransactionHandler struct {
	db *gorm.DB
}

func NewTransactionHandler(db *gorm.DB) *TransactionHandler {
	return &TransactionHandler{db: db}
}

func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	var transactions []models.BudgetTransaction

	// Busca transacciones cuyo status sea "Pendiente"
	if err := h.db.Preload("Budget").Preload("Budget.Client").Where("status = ?", "Pendiente").Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar transacciones", "detalles": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}

type AproveTransaction struct {
	TransactionID string `json:"transaction_id"`
}

func (h *TransactionHandler) AproveTransaction(c *gin.Context) {
	var req AproveTransaction

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Solicitud inválida", "detalles": err.Error()})
		return
	}

	// Buscamos la transacción junto con su presupuesto asociado
	var transaction models.BudgetTransaction
	if err := h.db.Preload("Budget").First(&transaction, "id = ?", req.TransactionID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar transacción", "detalles": err.Error()})
		return
	}

	// Solo se permite actualizar transacciones que estén en estado "Pendiente"
	if transaction.Status != "Pendiente" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transacción no está pendiente", "mensaje": "La transacción ya fue aprobada o rechazada"})
		return
	}

	// Actualizamos el balance del presupuesto asociado usando una consulta SQL
	sqlBudget := "UPDATE budgets SET balance = balance + ? WHERE id = ?"
	resultBudget := h.db.Exec(sqlBudget, transaction.Amount, transaction.BudgetID)
	if resultBudget.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar presupuesto", "detalles": resultBudget.Error.Error()})
		return
	}
	if resultBudget.RowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se encontró presupuesto para actualizar"})
		return
	}

	// Actualizamos el estado de la transacción a "Aprobada" usando una consulta SQL
	sqlTransaction := "UPDATE budget_transactions SET status = ? WHERE id = ?"
	resultTransaction := h.db.Exec(sqlTransaction, "Aprobada", transaction.ID)
	if resultTransaction.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar transacción", "detalles": resultTransaction.Error.Error()})
		return
	}
	if resultTransaction.RowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se encontró transacción para actualizar"})
		return
	}

	// Actualizamos el objeto local para reflejar el cambio y retornar la respuesta
	transaction.Status = "Aprobada"
	c.JSON(http.StatusOK, gin.H{"transaction": transaction})
}
