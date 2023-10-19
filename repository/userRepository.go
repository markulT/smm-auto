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
	var usersCollection = utils.DB.Collection("users")
	err := usersCollection.FindOne(context.TODO(), bson.M{"email":email}).Decode(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func SaveUser(user *models.User) error {
	var usersCollection = utils.DB.Collection("users")
	_, err := usersCollection.InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}
	return nil
}

func UpdateCustomerIDByEmail(email string, customerID string) error {
	var usersCollection = utils.DB.Collection("users")
	res := usersCollection.FindOneAndUpdate(context.TODO(), bson.M{"email":email}, bson.M{"$set":bson.M{"customerID":customerID}})
	return res.Err()
}

func UpdateUserSubscriptionID(email, subscriptionType, subscriptionID string) error {
	var usersCollection = utils.DB.Collection("users")
	res := usersCollection.FindOneAndUpdate(context.TODO(), bson.M{"email":email}, bson.M{"subscriptionID":subscriptionID, "subscriptionType":subscriptionType})
	return res.Err()
}

func GetUserSubLevelbyEmail(email string) (int, error) {
	var usersCollection = utils.DB.Collection("users")
	var user *models.User
	res := usersCollection.FindOne(context.TODO(),bson.M{"email":email})
	res.Decode(user)
	if res.Err() !=nil {
		return 0, res.Err()
	}
	return user.SubscriptionType, nil

}
