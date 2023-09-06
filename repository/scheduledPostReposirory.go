package repository

import (
	"context"
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



func GetScheduledPostRelations(offset int, limit int) *[]models.Post {
	var posts []models.Post
	postsCollection := utils.DB.Collection("posts")
	reqOptions:=options.Find()
	reqOptions.SetSkip(int64(offset))
	reqOptions.SetLimit(int64(limit))
	cur, err := postsCollection.Find(context.Background(), bson.M{"scheduled":bson.M{"$exists":true, "$ne":nil}}, reqOptions)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
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

func DeleteScheduledPostById(spId uuid.UUID) error {
	scheduledCollection := utils.DB.Collection("posts")
	_, err := scheduledCollection.DeleteOne(context.TODO(), bson.M{"_id":spId})
	if err != nil {
		return err
	}
	return nil
}

func SavePhoto(post *models.Post) (uuid.UUID, error) {
	postCollection := utils.DB.Collection("posts")
	res, err := postCollection.InsertOne(context.Background(), post)
	if err != nil {
		return uuid.UUID{}, err
	}
	id, err := uuid.Parse(res.InsertedID.(string))
	if err != nil {
		return uuid.UUID{},err
	}
	return id,nil
}

