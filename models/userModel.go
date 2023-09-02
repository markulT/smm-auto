package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	//ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Email    string `gorm:"unique"`
	Password string
}
