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
	//r.POST("/post", controllers.PostCreate)
	//r.GET("/post", controllers.PostIndex)
	//r.GET("/post/:id", controllers.PostShow)
	//r.PATCH("/post", controllers.PostUpdateWithReflect)
	//r.DELETE("/post/:id", controllers.PostDelete)
	controllers.SetupAuthRoutes(r)
	r.Run()
}
