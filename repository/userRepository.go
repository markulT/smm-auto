package repository

import (
	"context"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/bson"
	"golearn/models"
	"golearn/utils"
)

//type UserRepository interface {
//	SetUsersDeviceToken(userID uuid.UUID, token string) error
//	GetUserByEmail(email string) (models.User, error)
//	SaveUser(user *models.User) error
//	UpdateCustomerIDByEmail(email string, customerID string) error
//	UpdateUserSubscriptionID(email, subscriptionType, subscriptionID string) error
//	GetUserSubLevelbyEmail(email string) (int, error)
//	AddChannelToUser(userID,chID uuid.UUID) error
//}
//
//type userRepoImpl struct {}
//
//func NewUserRepo() UserRepository {
//	return &userRepoImpl{}
//}
//
//func (ur *userRepoImpl) AddChannelToUser(userID,chID uuid.UUID) error {
//	return nil
//}
//
//func (ur *userRepoImpl) SetUsersDeviceToken(userID uuid.UUID, token string) error {
//	var usersCollection = utils.DB.Collection("users")
//	err := usersCollection.FindOneAndUpdate(context.Background(),bson.M{"_id":userID}, bson.M{"deviceToken":token})
//	if err.Err()!=nil{
//		return err.Err()
//	}
//	return nil
//}
//
//func (ur *userRepoImpl) GetUserByEmail(email string) (models.User, error)  {
//	var user models.User
//	var usersCollection = utils.DB.Collection("users")
//	err := usersCollection.FindOne(context.TODO(), bson.M{"email":email}).Decode(&user)
//	if err != nil {
//		return user, err
//	}
//	return user, nil
//}
//
//func (ur *userRepoImpl) SaveUser(user *models.User) error {
//	var usersCollection = utils.DB.Collection("users")
//	_, err := usersCollection.InsertOne(context.TODO(), user)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//func (ur *userRepoImpl) UpdateCustomerIDByEmail(email string, customerID string) error {
//	var usersCollection = utils.DB.Collection("users")
//	res := usersCollection.FindOneAndUpdate(context.TODO(), bson.M{"email":email}, bson.M{"$set":bson.M{"customerID":customerID}})
//	return res.Err()
//}
//
//func (ur *userRepoImpl) UpdateUserSubscriptionID(email, subscriptionType, subscriptionID string) error {
//	var usersCollection = utils.DB.Collection("users")
//	res := usersCollection.FindOneAndUpdate(context.TODO(), bson.M{"email":email}, bson.M{"subscriptionID":subscriptionID, "subscriptionType":subscriptionType})
//	return res.Err()
//}
//
//func (ur *userRepoImpl) GetUserSubLevelbyEmail(email string) (int, error) {
//	var usersCollection = utils.DB.Collection("users")
//	var user *models.User
//	res := usersCollection.FindOne(context.TODO(),bson.M{"email":email})
//	res.Decode(user)
//	if res.Err() !=nil {
//		return 0, res.Err()
//	}
//	return user.SubscriptionType, nil
//
//}

//

func AddChannelToUser(c context.Context,userID,chID uuid.UUID) error {
	var usersCollection = utils.DB.Collection("users")
	_, err := usersCollection.UpdateOne(c, bson.M{"_id": userID}, bson.M{"$push": bson.M{"channelList": chID}})
	if err != nil {
		return err
	}
	return nil
}

func RemoveChannelFromUser(userID, chID uuid.UUID) error {
	var usersCollection = utils.DB.Collection("users")
	_, err := usersCollection.UpdateOne(context.Background(), bson.M{"_id": userID}, bson.M{"$pull": bson.M{"channelList": chID}})
	if err != nil {
		return err
	}
	return nil
}

func SetUsersDeviceToken(userID uuid.UUID, token string) error {
	var usersCollection = utils.DB.Collection("users")
	err := usersCollection.FindOneAndUpdate(context.Background(),bson.M{"_id":userID}, bson.M{"deviceToken":token})
	if err.Err()!=nil{
		return err.Err()
	}
	return nil
}



func GetUserByEmail(email string) (models.User, error)  {
	var user models.User
	var usersCollection = utils.DB.Collection("users")
	res := usersCollection.FindOne(context.Background(), bson.M{"email": email})

	err := res.Decode(&user)
	if err != nil {

		return models.User{}, err
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
