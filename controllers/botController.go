package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golearn/api/telegram"
)

func SetupBotRoutes(r *gin.Engine) {
	botGroup := r.Group("/bot")
	botGroup.GET("/getMe", getMeHandler)
	botGroup.POST("/sendMessage", sendMessageHandler)
}



func getMeHandler(c *gin.Context) {
	telegram.GetMe()
	c.JSON(200, gin.H{"message": "success"})
}
func sendMessageHandler(c *gin.Context)  {
	var body struct{
		Text string `json:"text"`
	}
	err := c.Bind(&body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(body.Text)
	telegram.SendMessage(body.Text)
	c.JSON(200, gin.H{"message":"successfully sent a message"})
}
