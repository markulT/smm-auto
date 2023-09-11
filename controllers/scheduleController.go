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
	"golearn/utils/videoCompress"
	"os"
	"path/filepath"
	"time"
)

func SetupScheduleRoutes(r *gin.Engine)  {
	scheduleGroup := r.Group("/schedule")
	scheduleGroup.POST("/message", scheduleMessageHandler)
	scheduleGroup.POST("/photo", schedulePhotoHandler)
	scheduleGroup.POST("/mediaGroup", scheduleMediaGroupHandler)
	scheduleGroup.POST("/video", scheduleVideoHandler)
}


func scheduleMessageHandler(c *gin.Context) {
	var body struct {
		Text string `json:"text"`
		Chat string `json:"chat"`
		Time string `json:"time"`
		Timezone string `json:"timezone"`
	}
	jsonHelper.BindWithException(&body, c)
	postId, _ := uuid.NewRandom()
	userId, _ := uuid.NewRandom()
	parsedTime, _ := time.Parse("2006 01-02 15:04 -0700 UTC", body.Time)
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
	scheduledTime := multipart.Value["scheduled"]
	parsedTime, _ := time.Parse("2006 01-02 15:04 -0700 UTC", scheduledTime[0])
	file, err := files[0].Open()
	defer file.Close()
	if err != nil {
		c.JSON(400, gin.H{"error":err})
		c.Abort()
		return
	}
	postID, err := uuid.NewRandom()
	if err != nil {
		c.JSON(400, gin.H{"error":err})
		c.Abort()
		return
	}
	post := models.Post{
		ID: postID,
		Text:        caption[0],
		ChannelName: channelName[0],
		Type:        "photo",
		Files:       []uuid.UUID{postID},
		Scheduled:   parsedTime,
	}
	err = mongoRepository.SaveScheduledPost(&post)
	if err != nil {
		c.JSON(400, gin.H{"error":err})
		c.Abort()
		return
	}
	err = s3.LoadImage(context.Background(), postID.String(), &file)
	if err != nil {
		c.JSON(400, gin.H{"error":err})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{"message":"success"})
}

func scheduleMediaGroupHandler(c *gin.Context)  {
	multipart, _ := c.MultipartForm()
	files := multipart.File["photo"]
	caption := multipart.Value["caption"]
	channelName := multipart.Value["channelName"]
	scheduledTime := multipart.Value["scheduled"]
	parsedTime, _ := time.Parse("2006 01-02 15:04 -0700 UTC", scheduledTime[0])
	postID, err := uuid.NewRandom()
	if err != nil {
		c.JSON(400, gin.H{"error":err})
		c.Abort()
		return
	}
	post := models.Post{
		ID: postID,
		Text:        caption[0],
		ChannelName: channelName[0],
		Type:        "mediaGroup",
		Files:       []uuid.UUID{},
		Scheduled:   parsedTime,
	}
	err = mongoRepository.SaveScheduledPost(&post)
	if err != nil {
		c.JSON(400, gin.H{"error":err})
		c.Abort()
		return
	}
	var fileIDList []uuid.UUID
	for _, file := range files{
		of, err := file.Open()
		if err != nil {
			c.JSON(400, gin.H{"error":err})
			c.Abort()
			return
		}

		fileID, err := uuid.NewRandom()
		if err != nil {
			c.JSON(400, gin.H{"error":err})
			c.Abort()
			return
		}
		err = s3.LoadMedia(context.Background(), fileID.String(),&of)
		if err != nil {
			c.JSON(400, gin.H{"error":err})
			c.Abort()
			return
		}
		fileIDList = append(fileIDList, fileID)
		err = of.Close()
		if err != nil {
			c.JSON(400, gin.H{"error":err})
			c.Abort()
			return
		}
	}
	err = mongoRepository.UpdateFilesList(postID, fileIDList)
	if err != nil {
		c.JSON(400, gin.H{"error":err})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{})
}

func scheduleVideoHandler(c *gin.Context)  {
	multipart, _ := c.MultipartForm()
	files := multipart.File["video"]
	caption := multipart.Value["caption"]
	channelName := multipart.Value["channelName"]
	scheduledTime := multipart.Value["scheduled"]

	parsedTime, _ := time.Parse("2006 01-02 15:04 -0700 UTC", scheduledTime[0])
	postID, err := uuid.NewRandom()
	if err != nil {
		c.JSON(400, gin.H{"error":err})
		c.Abort()
		return
	}
	post := models.Post{
		ID: postID,
		Text:        caption[0],
		ChannelName: channelName[0],
		Type:        "video",
		Files:       []uuid.UUID{},
		Scheduled:   parsedTime,
	}
	err = mongoRepository.SaveScheduledPost(&post)
	if err != nil {
		c.JSON(400, gin.H{"error":err})
		c.Abort()
		return
	}
	var fileIDList []uuid.UUID
	if files[0].Size > 48 * 1024 * 1024 {
		randomName, _ := uuid.NewRandom()
		compressedFileInfo, err := videoCompress.CompressFileToSize(files[0], randomName.String(), int64(48 * 1024 * 1024))
		if err != nil {
			c.JSON(400, gin.H{"error":err})
			c.Abort()
			return
		}
		of, err := os.Open(filepath.Join("C:/", compressedFileInfo.CompressedFilename))
		if err != nil {
			c.JSON(400, gin.H{"error":err})
			c.Abort()
			return
		}
		//_, err = telegram.SendVideo(of, caption[0], channelName[0], compressedFileInfo.CompressedFilename)
		//if err != nil {
		//	c.JSON(400, gin.H{"error":err})
		//	c.Abort()
		//	return
		//}

		fileID, err := uuid.NewRandom()
		err = s3.LoadVideo(fileID.String(), of)
		if err != nil {
			c.JSON(400, gin.H{"error":err})
			c.Abort()
			return
		}

		of.Close()
		err = videoCompress.CleanupCompressedFile(filepath.Join("C:/", compressedFileInfo.CompressedFilename))
		if err != nil {
			c.JSON(400, gin.H{"error":err})
			c.Abort()
			return
		}
		fileIDList = append(fileIDList, fileID)
	} else {
		of, err := files[0].Open()
		if err != nil {
			c.JSON(400, gin.H{"error":err})
			c.Abort()
			return
		}
		defer of.Close()
		fileID, err := uuid.NewRandom()
		err = s3.LoadVideoMultipart(fileID.String(), &of)
		if err != nil {
			c.JSON(400, gin.H{"error":err})
			c.Abort()
			return
		}
		fileIDList = append(fileIDList, fileID)
	}

	err = mongoRepository.UpdateFilesList(postID, fileIDList)
	if err != nil {
		c.JSON(400, gin.H{"error":err})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{})
}