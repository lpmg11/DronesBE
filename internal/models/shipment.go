package models

type Shipment struct {
	BaseModel
	DeliveryLocation  string   `gorm:"not null" json:"delivery_location"`
	DeliveryLatitude  float64  `gorm:"not null" json:"latitude"`
	DeliveryLongitude float64  `gorm:"not null" json:"longitude"`
	ShipmentCost      int      `gorm:"not null" json:"shipment_cost"`
	ProductID         uint     `gorm:"not null" json:"product_id"`
	Product           *Product `gorm:"foreignKey:ProductID;references:ID" json:"product"`
	DroneID           uint     `gorm:"not null" json:"drone_id"`
	Drone             *Drone   `gorm:"foreignKey:DroneID;references:ID" json:"drone"`
}
