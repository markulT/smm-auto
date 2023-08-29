package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golearn/api/telegram"
	"golearn/models"
	"golearn/utils"
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
func sendMessageHandler(c *gin.Context) {
	var body struct {
		Text        string `json:"text"`
		Scheduled   string `json:"time"`
		TimeZone    string `json:"timeZone"`
		ChannelName string `json:"channelName"`
		Username    string `json:"username"`
	}
	err := c.Bind(&body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(body.Text)
	fmt.Println(body.Scheduled)
	fmt.Println(body.TimeZone)
	fmt.Println(body.ChannelName)
	fmt.Println(body.Username)

	newPost := models.Post{
		Text:        body.Text,
		Scheduled:   body.Scheduled,
		Username:    body.Username,
		ChannelName: body.ChannelName,
		TimeZone:    body.TimeZone,
		Status:      "scheduled",
	}

	if err := utils.DB.Create(&newPost).Error; err != nil {
		c.JSON(400, gin.H{"message": "Error scheduling post"})
		return
	}

	c.JSON(200, gin.H{"message": "Message has been scheduled successfully"})

}
