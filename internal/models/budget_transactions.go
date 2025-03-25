package models

import "github.com/google/uuid"

type BudgetTransaction struct {
	BaseModel
	BudgetID         uuid.UUID `gorm:"type:uuid;not null" json:"budget_id"`
	Budget           *Budget   `gorm:"foreignKey:BudgetID;references:ID" json:"budget"`
	Amount           float64   `gorm:"not null" json:"amount"`
	Description      string    `gorm:"not null" json:"description"`
	ConfirmationCode string    `gorm:"not null" json:"confirmation_code"`
	Status           string    `gorm:"not null" json:"status"`
}
