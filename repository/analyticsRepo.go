package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"golearn/models"
	"golearn/utils"
	"time"
)

type AnalyticsRepo interface {
	GetManyByChannelIDAndDate(c context.Context,chID string, dateStart time.Time, dateEnd time.Time) ([]*models.AnalyticsUnit, error)
	SaveChannelAvgViews(c context.Context, channelName string, v float64) error
	GetAllOriginalChannels(c context.Context) ([]*models.Channel, error)
}

type defaultAnalyticsRepo struct {}

func NewAnalyticsRepo() AnalyticsRepo {
	return &defaultAnalyticsRepo{}
}

func (ar *defaultAnalyticsRepo) GetAllOriginalChannels(c context.Context) ([]*models.Channel, error) {

	analyticsCollection := utils.DB.Collection("analytics")

	//distinct, err := analyticsCollection.Distinct(context.Background(), "name", bson.M{})
	pipeline := []bson.M{
		{
			"$group": bson.M{
				"_id":     "$title",
				"doc":     bson.M{"$first": "$$ROOT"},
				"count":   bson.M{"$sum": 1},
				"uniqueId": bson.M{"$first": "$_id"},
			},
		},
		{
			"$replaceRoot": bson.M{"newRoot": "$doc"},
		},
	}

	channelsCursor, err := analyticsCollection.Aggregate(c, pipeline)
	if err != nil {
		return nil, err
	}
	defer channelsCursor.Close(context.Background())

	var result []*models.Channel
	for channelsCursor.Next(c) {
		var channel models.Channel
		if err:=channelsCursor.Decode(&channel);err!=nil {
			return nil, err
		}
		result = append(result, &channel)
	}

	return result, nil
}

func (ar *defaultAnalyticsRepo) SaveChannelAvgViews(c context.Context, channelName string, v float64) error {

	analyticsCollection := utils.DB.Collection("analytics")

	analyticsUnit := models.AnalyticsUnit{
		ChannelName: channelName,
		Value:       nil,
		Date:        time.Now(),
		Type:        "avgViews",
	}

	_, err := analyticsCollection.InsertOne(c, analyticsUnit)
	if err != nil {
		return err
	}

	return nil
}

func (ar *defaultAnalyticsRepo) GetManyByChannelIDAndDate(c context.Context,channelName string, dateStart time.Time, dateEnd time.Time) ([]*models.AnalyticsUnit, error) {
	fmt.Println("pidar")
	var analyticsUnitArray []*models.AnalyticsUnit
	analyticsCollection := utils.DB.Collection("analytics")
	fmt.Println("huy 1")
	fmt.Println(dateStart.Hour())
	fmt.Println("huy 2")
	curs, err := analyticsCollection.Find(c,bson.M{"channelName":channelName, "date":bson.M{
		"$gte": time.Date(dateStart.Year(), dateStart.Month(), dateStart.Day(), dateStart.Hour(), 0, 0,0 ,time.UTC),
		"$lt":dateEnd,
	}})
	fmt.Println("gowno")
	defer curs.Close(c)
	if err != nil {
		return nil, err
	}
	if curs.Err() != nil {
		return nil, curs.Err()
	}

	for curs.Next(c) {
		var analyticsUnit models.AnalyticsUnit
		if err:=curs.Decode(&analyticsUnit);err!=nil {
			return nil, err
		}
		fmt.Println("nigga")
		fmt.Println(analyticsUnitArray)
		analyticsUnitArray = append(analyticsUnitArray, &analyticsUnit)
	}
	fmt.Println(len(analyticsUnitArray))
	return analyticsUnitArray, nil
}

