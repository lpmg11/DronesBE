package models

type DroneModel struct{
	BaseModel
	ChargeKM    float64    `gorm:"not null" json:"charge_km"`
	Speed       float64    `gorm:"not null" json:"speed"`
}