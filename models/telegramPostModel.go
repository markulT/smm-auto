package models

import (
	"github.com/google/uuid"
	"time"
)


type Post struct {
	ID        	uuid.UUID `bson:"_id"`
	Text        string `bson:"text"`
	ChannelName string `bson:"channelName"`
	Type        string `bson:"type"`
	UserID		uuid.UUID `bson:"userId"`
	Files 		[]uuid.UUID `bson:"files"`
	Scheduled 	time.Time `bson:"scheduled,omitempty"`
}

//func (p Post) Value() (driver.Value, error) {
//	return p.ID, nil
//}
//
//func (p *Post) Scan(value interface{}) error {
//	p.ID = uint(value.(int64))
//	return nil
//}

