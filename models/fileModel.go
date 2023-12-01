package models

import (
	"github.com/google/uuid"
)

type File struct {
	ID       	uuid.UUID `bson:"_id" json:"id"`
	BucketName  string `bson:"bucketName" json:"bucketName"`
	Type 		string `bson:"type" json:"type"`
	PostID		uuid.UUID `bson:"postID,omitempty" json:"postId"`
	Filename 	string `bson:"filename" json:"filename"`
}
