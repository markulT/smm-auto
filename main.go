package main

import (
	"context"
	"github.com/gin-gonic/gin"
	_ "github.com/swaggo/files"
	swaggerFiles "github.com/swaggo/files"
	_ "github.com/swaggo/gin-swagger"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golearn/controllers"
	_ "golearn/docs"
	"golearn/utils"
	"golearn/utils/archiveCleaner"
	"golearn/utils/s3"
	"golearn/utils/scheduler"
	"os"
	"runtime"
	"runtime/pprof"
	"testing"
)

func TestMemoryLeak(t *testing.T) {
	runtime.MemProfileRate = 1

	// Your test logic here

	f, err := os.Create("memprofile")
	if err != nil {
		t.Fatal(err)
	}
	pprof.WriteHeapProfile(f)
	f.Close()
}

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
	controllers.SetupChannelRoutes(r)

	schedulerTask := &scheduler.SchedulerTask{}
	schedulerTask.FcmClient = firebaseMessagingClient

	go schedulerTask.FetchAndProcessPosts()
	go archiveCleaner.RunArchiveCleaner()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run()
}
