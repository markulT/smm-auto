package models

import "github.com/google/uuid"

type User struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Email string `gorm:"unique"`
	Password string
}
