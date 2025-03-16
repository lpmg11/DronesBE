package models

type User struct {
	BaseModel
	Username string `gorm:"uniqueIndex;not null" json:"username"`
	Password string `gorm:"not null" json:"-"`
	Role     string `gorm:"default:user;not null" json:"role"`
}
