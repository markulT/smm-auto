package models

type Plan struct {
	StripePlanID string `bson:"StripePlanID"`
	Level int `bson:"level"`
}
