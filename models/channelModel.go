package models

import "github.com/google/uuid"

type Channel struct {
	ID        uuid.UUID `bson:"_id" json:"id"`
	Name string `bson:"name" json:"name"`
	AssignedBotToken string `json:"assignedBotToken" bson:"assignedBotToken"`
	UserID uuid.UUID `bson:"userId" json:"userId"`
	ChannelHash string `bson:"channelHash,omitempty" json:"channelHash"`
}
