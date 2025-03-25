package models

type DroneModel struct {
	BaseModel
	Name       string  `gorm:"uniqueIndex;not null" json:"name"`
	ChargeKM   float64 `gorm:"not null" json:"charge_km"`
	Speed      float64 `gorm:"not null" json:"speed"`
	PricePerKm float64 `json:"price_per_km"`
}
