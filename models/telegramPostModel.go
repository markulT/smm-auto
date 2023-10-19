package models

import (
	"github.com/google/uuid"
	"time"
)


type Post struct {
	ID        	uuid.UUID `bson:"_id" json:"id"`
	Text        string `bson:"text" json:"content"`
	Title        string `bson:"title" json:"title"`
	ChannelName string `bson:"channelName" json:"chat"`
	Type        string `bson:"type" json:"type"`
	UserID		uuid.UUID `bson:"userId" bson:"userId"`
	Files 		[]uuid.UUID `bson:"files" bson:"files"`
	Scheduled 	time.Time `bson:"scheduled,omitempty" json:"scheduled"`
}
//a

type PostFile struct {
	ID uuid.UUID `bson:"file_id"`
	Type string `bson:"file_type"`
}

//func (p Post) Value() (driver.Value, error) {
//	return p.ID, nil
//}
//
//func (p *Post) Scan(value interface{}) error {
//	p.ID = uint(value.(int64))
//	return nil
//}

