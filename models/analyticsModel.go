package models

import (
	"github.com/google/uuid"
	"time"
)

type AnalyticsTypes string

const (
	AvgViews AnalyticsTypes = "avgViews"
)

type AnalyticsUnit struct {
	ID uuid.UUID `bson:"_id" json:"id"`
	ChannelName string `json:"channelName" bson:"channelName"`
	Value interface{} `json:"value" bson:"value"`
	Date time.Time `json:"date" bson:"date"`
	Type AnalyticsTypes `json:"type" bson:"type"`
}

func NewAnalyticsUnit(channelName string, value float64, at AnalyticsTypes) AnalyticsUnit {
	return AnalyticsUnit{
		ChannelName: channelName,
		Value:       value,
		Date:        time.Now(),
		Type:        at,
	}
}

