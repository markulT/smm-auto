package models

import (
	"github.com/google/uuid"
	"time"
)

type Role string

const (
	RoleAdmin     Role = "admin"
	RoleUser      Role = "user"
	RoleModerator Role = "moderator"
	RoleSubUser   Role = "sub_user"
)

type User struct {

	ID               uuid.UUID `bson:"_id" json:"id"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Email            string
	Password         string
	SubscriptionID   string `bson:"subscriptionID"`
	SubscriptionType int    `bson:"subscriptionType"`
	Role             Role
	CustomerID       string `bson:"customerID"`
	DeviceToken      string `bson:"deviceToken" json:"deviceToken"`

}
