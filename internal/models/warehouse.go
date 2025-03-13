package models

type Warehouse struct {
	BaseModel
	Name      string   `gorm:"uniqueIndex;not null" json:"name"`
	Latitude  float64  `gorm:"not null" json:"latitude"`
	Longitude float64  `gorm:"not null" json:"longitude"`
	Drones    []*Drone `gorm:"foreignKey:WarehouseID;references:ID" json:"drones"`
}
