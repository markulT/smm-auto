package controllers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golearn/models"
	mongoRepository "golearn/repository"
	"golearn/utils/jsonHelper"
	"golearn/utils/s3"
	"time"
)

func SetupScheduleRoutes(r *gin.Engine)  {
	scheduleGroup := r.Group("/schedule")
	scheduleGroup.POST("/message", scheduleMessageHandler)
	scheduleGroup.POST("/photo", schedulePhotoHandler)
}


func scheduleMessageHandler(c *gin.Context) {
	var body struct {
		Text string `json:"text"`
		Chat string `json:"chat"`
		Time string `json:"time"`
	}
	jsonHelper.BindWithException(&body, c)
	postId, _ := uuid.NewRandom()
	userId, _ := uuid.NewRandom()
	parsedTime, _ := time.Parse("2006 01-02 15:04", body.Time)
	post:=models.Post{
		Text:        body.Text,
		ChannelName: body.Chat,
		Type:        "message",
		UserID:      userId,
		Files:    	 nil,
		ID: 		 postId,
		Scheduled: parsedTime,
	}
	fmt.Println(postId)
	err := mongoRepository.SavePostWithId(&post, postId)
	if err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{"message":"Error saving post"})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{"message":"message"})
}

func schedulePhotoHandler(c *gin.Context)  {
	multipart, _ := c.MultipartForm()

	files := multipart.File["photo"]
	caption := multipart.Value["caption"]
	channelName := multipart.Value["channelName"]

	file, err := files[0].Open()
	defer file.Close()
	if err != nil {
		c.JSON(400, gin.H{"error":err})
		c.Abort()
		return
	}

	post := models.Post{
		Text:        caption[0],
		ChannelName: channelName[0],
		Type:        "photo",
		Files:       nil,
		Scheduled:   time.Time{},
	}
	savedId, err := mongoRepository.SavePhoto(&post)
	if err != nil {
		c.JSON(400, gin.H{"error":err})
		c.Abort()
		return
	}
	err = s3.LoadImage(context.Background(), savedId.String(), &file)
	if err != nil {
		c.JSON(400, gin.H{"error":err})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{"message":"success"})
}

func scheduleMediaGroupHandler(c *gin.Context)  {

}