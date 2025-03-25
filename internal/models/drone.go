package models

import "github.com/google/uuid"

type Drone struct {
	BaseModel
	Name        string      `gorm:"uniqueIndex;not null" json:"name"`
	WarehouseID uuid.UUID   `gorm:"not null" json:"warehouse_id"`
	Warehouse   *Warehouse  `gorm:"foreignKey:WarehouseID;references:ID" json:"warehouse"`
	ModelId     uuid.UUID   `gorm:"not null" json:"model_id"`
	Model       *DroneModel `gorm:"foreingKey:" json:"model"`
}
