package repository

import (
	"context"
	"github.com/google/uuid"
	"golearn/models"
	"golearn/utils"
)

func SavePost(post *models.Post) error {
	postCollection := utils.DB.Collection("posts")
	_, err := postCollection.InsertOne(context.TODO(), post)
	if err != nil {
		return err
	}
	return nil
}

func SavePostWithId(post *models.Post, id uuid.UUID) error {
	postCollection := utils.DB.Collection("posts")
	_, err := postCollection.InsertOne(context.TODO(), post)
	if err != nil {
		return err
	}
	return nil
}