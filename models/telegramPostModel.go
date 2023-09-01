package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	//ID          uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Text        string
	Scheduled   string
	TimeZone    string
	ChannelName string
	Username    string
	Status      string
	Type        string
}
