package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"golearn/controllers"
	"golearn/utils"
	"golearn/utils/archiveCleaner"
	"golearn/utils/s3"
	"golearn/utils/scheduler"
)

func init() {
	utils.LoadEnvVariables()
	utils.ConnectToDb()
	s3.ConnectToMinio()
	utils.StripeInit()
}

// @title SMM-Auto API
// @version 1.0
// @description API server for SMM-auto application

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
func main() {

	firebaseApp, err := utils.FirebaseInit()
	if err != nil {
		panic(err)
	}
	firebaseMessagingClient, err := firebaseApp.Messaging(context.Background())
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	controllers.SetupAuthRoutes(r)
	controllers.SetupBotRoutes(r)
	controllers.SetupScheduleRoutes(r)
	controllers.SetupPaymentRoutes(r)
	controllers.SetupArchiveRoutes(r)

	schedulerTask := &scheduler.SchedulerTask{}
	schedulerTask.FcmClient = firebaseMessagingClient

	go schedulerTask.FetchAndProcessPosts()
	go archiveCleaner.RunArchiveCleaner()

	r.Run()
}
