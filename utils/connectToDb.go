package utils

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

//var DB *gorm.DB
//
//func ConnectToDb() {
//	var err error
//	dsn := os.Getenv("DB_URL")
//	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
//
//	if err != nil {
//		log.Fatal("Error while connecting to database")
//	}
//}
//

var DB *mongo.Database

func ConnectToDb() {

	mongoUri := os.Getenv("DB_URL")
	mongoOptions:=options.Client().ApplyURI(mongoUri)
	client, err := mongo.Connect(context.TODO(), mongoOptions)
	if err != nil {
		log.Println("Error while connecting to database at /utils/connectToDb.go")
		log.Fatal(err)
	}
	DB = client.Database("smm-auto")

}