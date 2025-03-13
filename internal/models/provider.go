package models

import "github.com/google/uuid"

type Provider struct {
	BaseModel
	Name      string    `gorm:"uniqueIndex;not null" json:"name"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Latitude  float64   `gorm:"not null" json:"latitude"`
	Longitude float64   `gorm:"not null" json:"longitude"`
	User      *User     `gorm:"foreingKey:UserID;references:ID" json:"user"`
}
