package main

import (
	"github.com/gin-gonic/gin"
	"golearn/controllers"
	"golearn/utils"
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
	r := gin.Default()

	controllers.SetupAuthRoutes(r)
	controllers.SetupBotRoutes(r)
	controllers.SetupScheduleRoutes(r)
	controllers.SetupPaymentRoutes(r)

	go scheduler.FetchAndProcessPosts()

	r.Run()
}
