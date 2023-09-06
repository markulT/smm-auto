package models

import (
	"time"
)

type User struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	Email    string
	Password string
	ChannelList []byte
	SubscriptionID string
	SubscriptionType string
}

