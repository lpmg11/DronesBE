package models

import "github.com/google/uuid"

type Product struct {
	BaseModel
	Name        string    `gorm:"uniqueIndex;not null" json:"name"`
	Description string    `gorm:"not null" json:"description"`
	Price       int       `gorm:"not null" json:"price"`
	ProviderID  uuid.UUID `gorm:"type:uuid;not null" json:"provider_id"`
	Provider    *Provider `gorm:"foreignKey:ProviderID;references:ID" json:"provider"`
}
