package repository

import (
	"context"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golearn/models"
	"golearn/utils"
)


type ChannelRepository interface {
	AssignBotToken(c context.Context, botToken string, chID uuid.UUID) error
	SaveNewChannel(c context.Context, ch *models.Channel) error
	DeleteChannel(c context.Context, chID uuid.UUID) error
	FindByID(c context.Context, chID uuid.UUID) (*models.Channel, error)
	GetAllByUserID(c context.Context, userID uuid.UUID) ([]models.Channel, error)
}

type channelRepoImpl struct {}

func NewChannelRepo() ChannelRepository {
	return &channelRepoImpl{}
}

func (cr *channelRepoImpl) GetAllByUserID(c context.Context, userID uuid.UUID) ([]models.Channel, error) {
	var channelList []models.Channel

	channelCollection := utils.DB.Collection("channels")
	channelCursor,err := channelCollection.Find(c, bson.M{"userId":userID})
	defer channelCursor.Close(c)
	if err != nil {
		return nil, err
	}
	if channelCursor.Err()!=nil {
		return nil, err
	}

	for channelCursor.Next(c) {
		var channel models.Channel
		if err:=channelCursor.Decode(&channel);err!=nil {
			return nil, err
		}
		channelList = append(channelList, channel)
	}
	return channelList, nil
}

func (cr *channelRepoImpl) AssignBotToken(c context.Context, botToken string, chID uuid.UUID) error {
	channelCollection := utils.DB.Collection("channels")
	res := channelCollection.FindOneAndUpdate(c, bson.M{"_id":chID}, bson.M{"assignedBotToken":botToken})
	if res.Err() != nil {
		return res.Err()
	}
	return nil
}

func (cr *channelRepoImpl) SaveNewChannel(c context.Context, ch *models.Channel) error {
	session, err := utils.DB.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(c)

	_, err = session.WithTransaction(c, func(ctx mongo.SessionContext) (interface{}, error) {
		channelCollection := utils.DB.Collection("channels")
		_, err = channelCollection.InsertOne(ctx, ch)
		if err != nil {
			return nil, err
		}
		err = AddChannelToUser(ctx, ch.UserID, ch.ID)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (cr *channelRepoImpl) DeleteChannel(c context.Context, chID uuid.UUID) error {
	channelCollection := utils.DB.Collection("channels")
	_, err := channelCollection.DeleteOne(c, bson.M{"_id":chID})
	if err != nil {
		return err
	}
	return nil
}

func (cr *channelRepoImpl) FindByID(c context.Context, chID uuid.UUID) (*models.Channel, error) {
	var searchedChannel *models.Channel
	channelCollection := utils.DB.Collection("channels")
	res := channelCollection.FindOne(c, bson.M{"_id":chID})
	if res.Err() != nil {
		return nil, res.Err()
	}

	err := res.Decode(searchedChannel)
	if err != nil {
		return nil, err
	}

	return searchedChannel, nil
}