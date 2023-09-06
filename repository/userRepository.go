package repository

import (
	"context"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/bson"
	"golearn/models"
	"golearn/utils"
)



func GetUserByEmail(email string) (models.User, error)  {
	var user models.User
	usersCollection := utils.DB.Collection("users")
	err := usersCollection.FindOne(context.TODO(), bson.M{"email":email}).Decode(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}
func SaveUser(user *models.User) error {


	usersCollection := utils.DB.Collection("users")
	_, err := usersCollection.InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}
	return nil
}
