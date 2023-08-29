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
	botGroup.POST("/sendPhoto", sendPhotoHandler)
	botGroup.POST("/sendMediaGroup", sendMediaGroupHandler)
	botGroup.POST("/sendDice", sendDiceHandler)
	botGroup.POST("/test", sendMediaGroupLinks)
}

func sendDiceHandler(c *gin.Context)  {
	telegram.SendDice()
	c.JSON(200, gin.H{"a":"a"})
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
func sendPhotoHandler(c *gin.Context) {

	multipart, _ := c.MultipartForm()
	files := multipart.File["photo"]
	caption := multipart.Value["caption"]
	for _, file := range files {
		telegram.SendPhoto(file, caption[0])
	}
	c.JSON(200, gin.H{"aboba":"aboba"})
}
func sendMediaGroupHandler(c *gin.Context) {
	fmt.Println("in request")
	multipart, _ := c.MultipartForm()
	files := multipart.File["photo"]
	caption := multipart.Value["caption"]
	_, err := telegram.SendMediaGroup(files, caption[0])
	if err!= nil {
		c.JSON(500, gin.H{"error":err})
	}
	c.JSON(200, gin.H{"message":"success"})
}

func sendMediaGroupLinks(c *gin.Context) {
	_, _ = telegram.SendMediaGroupLinks([]string{"https://i.pinimg.com/originals/61/3b/f8/613bf893ab736ac25c6f6dde1bbacc4a.jpg", "https://sweetpeaskitchen.com/wp-content/uploads/2020/06/Choc-Rasp-Cheesecake-7-1-scaled.jpg"}, "text")
	c.JSON(200, gin.H{"message":"success"})
}