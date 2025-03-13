package models

import "github.com/google/uuid"

type Drone struct {
	BaseModel
	Name        string     `gorm:"uniqueIndex;not null" json:"name"`
	ChargeKM    float64    `gorm:"not null" json:"charge_km"`
	Speed       float64    `gorm:"not null" json:"speed"`
	WarehouseID uuid.UUID  `gorm:"not null" json:"warehouse_id"`
	Warehouse   *Warehouse `gorm:"foreignKey:WarehouseID;references:ID" json:"warehouse"`
}
