package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golearn/models"
	"golearn/utils"
	"log"
	"time"
)
type FlatScheduledRelations struct {
	SpId uuid.UUID `db:"sp_id"`
	SpTime time.Time `db:"time"`
	SpPostId uuid.UUID `db:"post_id"`
	PId uuid.UUID `db:"id"`
	PText        string `db:"text"`
	PChannelName string `db:"channel_name"`
	PType        string `db:"type"`
	PUserID		uuid.UUID `db:"user_id"`
	PFiles 		[]uuid.UUID `db:"files"`
	PScheduleID	int64 `db:"schedule_id"`
}



func GetScheduledPostRelations(c context.Context, offset int, limit int, archived bool) *[]models.Post {
	var posts []models.Post
	postsCollection := utils.DB.Collection("posts")
	reqOptions:=options.Find()
	reqOptions.SetSkip(int64(offset))
	reqOptions.SetLimit(int64(limit))
	cur, err := postsCollection.Find(c, bson.M{"scheduled":bson.M{"$exists":true, "$ne":nil}, "archived":archived}, reqOptions)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer cur.Close(c)
	for cur.Next(c) {
		var post models.Post
		if err:=cur.Decode(&post);err!=nil {
			log.Fatal(err)
		}
		posts = append(posts, post)
	}
	if err:=cur.Err();err!=nil {
		log.Fatal(err)
	}
	return &posts
}

func DeleteScheduledPostById(c context.Context,spId uuid.UUID) error {
	scheduledCollection := utils.DB.Collection("posts")
	_, err := scheduledCollection.DeleteOne(c, bson.M{"_id":spId})
	if err != nil {
		return err
	}
	return nil
}

func SaveScheduledPost(c context.Context,post *models.Post) error {
	postCollection := utils.DB.Collection("posts")
	_, err := postCollection.InsertOne(c, post)
	if err != nil {
		return err
	}
	return nil
}

func UpdateFilesList(c context.Context,pId uuid.UUID, files []uuid.UUID) error {
	postCollection := utils.DB.Collection("posts")
	_, err := postCollection.UpdateByID(c, pId, bson.M{"$set":bson.M{"files":files}})
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func ArchivizePost(c context.Context,pid uuid.UUID) error {
	postCollection := utils.DB.Collection("posts")
	_, err := postCollection.UpdateByID(c, pid, bson.M{"archived":true})
	if err != nil {
		return err
	}
	return nil
}
