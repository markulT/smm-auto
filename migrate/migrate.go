package main

import (
	"golearn/models"
	"golearn/utils"
)

func init() {
	utils.LoadEnvVariables()
	utils.ConnectToDb()
}

func main() {
	utils.DB.AutoMigrate(&models.User{})
	utils.DB.AutoMigrate(&models.Post{})
	utils.DB.AutoMigrate(&models.File{})
	utils.DB.Debug().AutoMigrate(&models.ScheduledPost{})
}
