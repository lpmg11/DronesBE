package models

import "github.com/google/uuid"

type Budget struct {
	BaseModel
	ClientID     uuid.UUID            `gorm:"type:uuid;not null" json:"client_id"`
	Client       *Client              `gorm:"foreignKey:ClientID;references:ID" json:"client"`
	Balance      float64              `gorm:"not null" json:"balance"`
	Transactions []*BudgetTransaction `gorm:"foreignKey:BudgetID;references:ID" json:"transactions"`
}
