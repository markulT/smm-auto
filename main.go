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
	"golearn/repository"
	"golearn/utils"
	"golearn/utils/analytics"
	"golearn/utils/archiveCleaner"
	"golearn/utils/payments"
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

	paymentRepo:=repository.NewPaymentRepo()
	paymentService:=payments.NewStripePaymentService(paymentRepo)
	chRepo := repository.NewChannelRepo()
	analyticsRepo := repository.NewAnalyticsRepo()
	analyticsService := analytics.NewAnalyticsService(analyticsRepo)

	controllers.SetupAuthRoutes(r)
	controllers.SetupAnalyticsRoutes(r, analyticsService, analyticsRepo)
	controllers.SetupArchiveRoutes(r)
	controllers.SetupBotRoutes(r)
	controllers.SetupScheduleRoutes(r)
	controllers.SetupPaymentRoutes(r, paymentService, paymentRepo)
	controllers.SetupChannelRoutes(r)


	schedulerTask := &scheduler.SchedulerTask{}
	schedulerTask.FcmClient = firebaseMessagingClient
	schedulerTask.ChRepo = chRepo



	analyticsTask := &analytics.AnalyticsTask{}
	analyticsTask.SetRepo(analytics.CreateBridgeRepo(analyticsRepo))

	go schedulerTask.FetchAndProcessPosts()
	go archiveCleaner.RunArchiveCleaner()
	//go analyticsTask.RunAnalyticsModule()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run()
}
