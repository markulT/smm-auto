package utils

import (
	"context"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func FirebaseInit() (*firebase.App, error) {
	opt:= option.WithCredentialsFile("/public/serviceAccountKey.json")
	app , err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}
	return app, nil
}
