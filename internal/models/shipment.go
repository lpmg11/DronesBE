package models

type Shipment struct {
	BaseModel
	ProductID uint     `gorm:"not null" json:"product_id"`
	Product   *Product `gorm:"foreignKey:ProductID;references:ID" json:"product"`
	DroneID   uint     `gorm:"not null" json:"drone_id"`
	Drone     *Drone   `gorm:"foreignKey:DroneID;references:ID" json:"drone"`
}
