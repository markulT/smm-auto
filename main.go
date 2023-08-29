package main

import (
	"github.com/gin-gonic/gin"
	"golearn/controllers"
	"golearn/utils"
)

func init() {
	utils.LoadEnvVariables()
	utils.ConnectToDb()
}



func main() {
	r := gin.Default()

	controllers.SetupAuthRoutes(r)
	controllers.SetupTelegramRoutes(r)
	controllers.SetupBotRoutes(r)
	r.Run()
}

