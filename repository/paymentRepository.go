package repository

import (
	"context"
	"github.com/stripe/stripe-go/v75"
	"go.mongodb.org/mongo-driver/bson"
	"golearn/models"
	"golearn/utils"
)



type PaymentRepo interface {
	SaveSubscription(c context.Context,sub models.Subscription) error
	DeleteSubscriptionByID(c context.Context, subID string) error
	FindSubscriptionByCustomerID(c context.Context, cID string) (*models.Subscription , error)
}

type defaultPaymentRepo struct {}

func NewPaymentRepo() PaymentRepo {
	return &defaultPaymentRepo{}
}

func (pr *defaultPaymentRepo) FindSubscriptionByCustomerID(c context.Context, cID string) (*models.Subscription , error) {
	var s stripe.Subscription
	paymentCollection := utils.DB.Collection("payments")
	res := paymentCollection.FindOne(c, bson.M{"customerId": cID})
	if res.Err() != nil {
		return nil, res.Err()
	}
	if err :=res.Decode(&s);err!=nil {
		return nil, err
	}
	subModel, err := models.NewSubscriptionFromStripe(&s)
	if err != nil {
		return nil, err
	}
	return subModel, nil
}

func (pr *defaultPaymentRepo) SaveSubscription(c context.Context,sub models.Subscription) error {

	paymentCollection := utils.DB.Collection("payments")

	_,err := paymentCollection.InsertOne(c, sub)
	if err != nil {
		return err
	}
	return nil
}

func (pr *defaultPaymentRepo) DeleteSubscriptionByID(c context.Context, subID string) error {
	paymentCollection := utils.DB.Collection("payments")
	_, err := paymentCollection.DeleteOne(c, bson.M{"_id":subID})
	if err != nil {
		return err
	}
	return nil
}
