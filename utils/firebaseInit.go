package utils

import (
	"context"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

//
//func FirebaseInit() (*firebase.App, error) {
//	opt:= option.WithCredentialsFile("./public/serviceAccountKey.json")
//	app , err := firebase.NewApp(context.Background(), nil, opt)
//	if err != nil {
//		return nil, err
//	}
//	return app, nil
//}

func FirebaseInit() (*firebase.App, error) {
	// Read the contents of the JSON file
	opt:= option.WithCredentialsFile("public/smm-auto-firebase-adminsdk-80m2e-b3964f7528.json")
	//opts := []option.ClientOption{option.WithCredentialsJSON([]byte(`{"apiKey":"AIzaSyDn5tALtyUVfxy4IVmGiL77-w47ewsh604", "projectId": "smm-auto"}`))}

	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}
	return app, nil
}
