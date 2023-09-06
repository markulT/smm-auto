package models

import (
	"github.com/google/uuid"
)

type File struct {
	ID       	uuid.UUID `bson:"_id"`
	BucketName  string `bson:"bucketName"`
	Type 		string `bson:"type"`
	PostID		uuid.UUID `bson:"postID,omitempty"`
}
