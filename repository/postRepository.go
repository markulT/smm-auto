package repository

import (
	"context"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"golearn/models"
	"golearn/utils"
	"sync"
	"time"
)

type PostRepository interface {
	SavePost(*models.Post) error
	SavePostWithId(*models.Post, uuid.UUID) error
	GetPostsByUserID(context.Context, uuid.UUID, *sync.WaitGroup, chan []models.Post)
	GetPostByID(context.Context, uuid.UUID) (models.Post, error)
	DeletePostByID(ctx context.Context, uuid2 uuid.UUID) bool
	GetPostByImageName(ctx context.Context, imageName uuid.UUID) (models.Post,error)
	GetPostsByDate(c context.Context, scheduled time.Time,userId uuid.UUID, wg *sync.WaitGroup, respch chan []models.Post)
	GetAllArchivedPostsByUserID(c context.Context, userID uuid.UUID) ([]models.Post, error)
}

type postRepositoryImpl struct {

}

func NewPostRepository() PostRepository {
	return &postRepositoryImpl{}
}

func (p *postRepositoryImpl) GetAllArchivedPostsByUserID(c context.Context, userID uuid.UUID) ([]models.Post, error) {
	postCollection := utils.DB.Collection("posts")
	postsCursor, err := postCollection.Find(c, bson.M{"userId":userID})
	if err != nil {
		return nil, err
	}
	defer postsCursor.Close(c)
	var posts []models.Post
	for postsCursor.Next(c) {
		var post models.Post
		if err := postsCursor.Decode(&post);err!=nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := postsCursor.Err();err!=nil {
		return nil, err
	}

	return posts, nil
}

func (p *postRepositoryImpl) GetPostsByDate(c context.Context, scheduled time.Time, userId uuid.UUID, wg *sync.WaitGroup, respch chan []models.Post )  {
	wg.Add(1)
	var posts []models.Post
	postCollection := utils.DB.Collection("posts")
	results, err := postCollection.Find(c, bson.M{
		"userId": userId,
		"date": bson.M{
			"$gte": time.Date(scheduled.Year(), scheduled.Month(), scheduled.Day(), 0, 0 ,0, 0, time.UTC),
			"$lte": time.Date(scheduled.Year(), scheduled.Month(), scheduled.Day(), 23, 59 ,59, 999999999, time.UTC),
		},
	})
	if err != nil {
		respch<-nil
		return
	}
	defer results.Close(c)
	for results.Next(c) {
		var post models.Post
		if err := results.Decode(&post);err!=nil {
			respch <- nil
			return
		}
		posts = append(posts, post)
	}

	if err := results.Err();err!=nil {
		respch<-nil
		return
	}
	respch <- posts
	wg.Done()
}

func (p *postRepositoryImpl) GetPostByImageName(c context.Context, imageName uuid.UUID) (models.Post,error) {
	var post models.Post
	postCollection := utils.DB.Collection("posts")
	res := postCollection.FindOne(c, bson.M{"files":bson.M{"$elemMatch":bson.M{"$eq":imageName}}})
	if err:=res.Decode(&post);err!=nil {
		return models.Post{},res.Err()
	}
	return post, nil
}

func (p *postRepositoryImpl) DeletePostByID(c context.Context, postID uuid.UUID) bool {
	postCollection := utils.DB.Collection("posts")
	_, err := postCollection.DeleteOne(c, bson.M{"_id": postID})
	if err != nil {
		return false
	}
	return true
}

func (p *postRepositoryImpl) GetPostByID(c context.Context, postID uuid.UUID) (models.Post, error) {
	postCollection := utils.DB.Collection("posts")
	var post models.Post
	err := postCollection.FindOne(c, bson.M{"_id":postID}).Decode(&post)
	if err != nil {
		return models.Post{},err
	}
	return post, nil
}

func (p *postRepositoryImpl) SavePost(post *models.Post) error {
	postCollection := utils.DB.Collection("posts")
	_, err := postCollection.InsertOne(context.TODO(), post)
	if err != nil {
		return err
	}
	return nil
}

func (p *postRepositoryImpl) SavePostWithId(post *models.Post, id uuid.UUID) error {
	postCollection := utils.DB.Collection("posts")
	_, err := postCollection.InsertOne(context.TODO(), post)
	if err != nil {
		return err
	}
	return nil
}

func (p *postRepositoryImpl) GetPostsByUserID(c context.Context,userID uuid.UUID, wg *sync.WaitGroup, respch chan []models.Post) {
	wg.Add(1)
	postCollection := utils.DB.Collection("posts")
	postsCursor, err := postCollection.Find(c, bson.M{"userId":userID})
	if err != nil {
		respch <- nil
		return
	}
	defer postsCursor.Close(c)
	var posts []models.Post
	for postsCursor.Next(c) {
		var post models.Post
		if err := postsCursor.Decode(&post);err!=nil {
			respch <- nil
			return
		}
		posts = append(posts, post)
	}

	if err := postsCursor.Err();err!=nil {
		respch<-nil
		return
	}
	respch<-posts
	wg.Done()
}

