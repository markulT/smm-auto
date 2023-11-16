package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"golearn/utils"
)

func WithTransaction(c context.Context, fn func(ctx mongo.SessionContext) (interface{}, error) ) error {
	session, err := utils.DB.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(c)

	_, err = session.WithTransaction(c, fn)
	if err != nil {
		return err
	}
	return nil
}
