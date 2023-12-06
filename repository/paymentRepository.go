package repository

import (
	"context"
	"golearn/models"
	"golearn/utils"
)



type PaymentRepo interface {
	SaveSubscription(c context.Context,sub models.Subscription) error
}

type defaultPaymentRepo struct {}

func NewPaymentRepo() PaymentRepo {
	return &defaultPaymentRepo{}
}

func (pr *defaultPaymentRepo) SaveSubscription(c context.Context,sub models.Subscription) error {

	paymentCollection := utils.DB.Collection("payments")

	_,err := paymentCollection.InsertOne(c, sub)
	if err != nil {
		return err
	}
	return nil
}
