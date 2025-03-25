package models

import "github.com/google/uuid"

type Client struct {
	BaseModel
	Name      string    `gorm:"uniqueIndex;not null" json:"name"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;unique" json:"user_id"`
	User      *User     `gorm:"foreingKey:UserID;references:ID" json:"user"`
	Budget    *Budget   `gorm:"foreignKey:ClientID;references:ID" json:"budget"`
	Latitude  float64   `gorm:"not null" json:"latitude"`
	Longitude float64   `gorm:"not null" json:"longitude"`
}
