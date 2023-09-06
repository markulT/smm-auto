package models

import (
	"github.com/google/uuid"
	"time"
)

type ScheduledPost struct {
	ID 			uuid.UUID `db:"sp_id"`
	PostID   	uuid.UUID `db:"sp_post_id"`
	Time		time.Time `db:"sp_post_id"`
}
type PostScheduleRelation struct {
	ScheduledPost ScheduledPost
	Post	Post
}