package models

type Drone struct {
	BaseModel
	Name     string  `gorm:"uniqueIndex;not null" json:"name"`
	ChargeKM float64 `gorm:"not null" json:"charge_km"`
	Speed    float64 `gorm:"not null" json:"speed"`
}
