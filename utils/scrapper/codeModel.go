package scrapper

import (
	"github.com/google/uuid"
	"time"
)

type Code struct {
	ID          uuid.UUID `bson:"_id" json:"id"`
	ChannelName string    `bson:"channelName" json:"channelName"`
	SubmitCode  string    `bson:"submitCode" json:"submitCode"`
	ExpireDate  time.Time `bson:"expireDate" json:"expireDate"`
}
