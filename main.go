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
}

func main() {
	r := gin.Default()

	controllers.SetupAuthRoutes(r)
	controllers.SetupTelegramRoutes(r)
	controllers.SetupBotRoutes(r)
	controllers.SetupScheduleRoutes(r)

	go scheduler.FetchAndProcessPosts()

	r.Run()
}
