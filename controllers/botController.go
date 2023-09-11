package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golearn/api/telegram"
	"golearn/models"
	"golearn/repository"
	"golearn/utils/videoCompress"
	"io"
	"os"
	"path/filepath"
)

func SetupBotRoutes(r *gin.Engine) {
	botGroup := r.Group("/bot")
	botGroup.GET("/getMe", getMeHandler)
	botGroup.POST("/sendMessage", sendMessageHandler)
	botGroup.POST("/sendPhoto", sendPhotoHandler)
	botGroup.POST("/sendMediaGroup", sendMediaGroupHandler)
	botGroup.POST("/sendDice", sendDiceHandler)
	botGroup.POST("/test", sendMediaGroupLinks)
	botGroup.POST("/sendAudio", sendAudioHandler)
	botGroup.POST("/sendVoice", sendVoiceHandler)
	botGroup.POST("/sendVideo", sendVideoHandler)
	botGroup.POST("/sendVideoNote", sendVideoNoteHandler)
	botGroup.POST("/sendLocation", sendLocationHandler)
	botGroup.POST("/sendVenue", sendVenueHandler)
	//botGroup.DELETE("/delete/:id", postDelete)
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

	newPost := models.Post{
		Text:        body.Text,
		ChannelName: body.ChannelName,
		Type:        "message",
	}

	if err := repository.SavePost(&newPost); err != nil {
		c.JSON(400, gin.H{"message": "Error scheduling post"})
		return
	}

	c.JSON(200, gin.H{"message": "Message has been scheduled successfully"})

}

func sendAudioHandler(c *gin.Context) {
	multipart, _ := c.MultipartForm()
	files := multipart.File["audio"]
	caption := multipart.Value["caption"]
	channelName := multipart.Value["channelName"]
	for _, file := range files {
		telegram.SendAudio(file, caption[0], channelName[0])
	}
	c.JSON(200, gin.H{"message": "success"})
}

func sendVoiceHandler(c *gin.Context) {
	multipart, _ := c.MultipartForm()
	files := multipart.File["voice"]
	caption := multipart.Value["caption"]
	channelName := multipart.Value["channelName"]
	for _, file := range files {
		telegram.SendVoice(file, caption[0], channelName[0])
	}
	c.JSON(200, gin.H{"message": "success"})
}

func sendVideoHandler(c *gin.Context) {
	fmt.Println("aboba")
	multipart, _ := c.MultipartForm()
	files := multipart.File["file"]
	caption := multipart.Value["caption"]
	channelName := multipart.Value["channelName"]
	//filename, err := uuid.NewRandom()
	//if err != nil {
	//	c.JSON(400, gin.H{"error":err})
	//	c.Abort()
	//	return
	//}
	fmt.Println(files[0].Size)
	if files[0].Size > 48 * 1024 * 1024 {
		fmt.Println("bigger while")
		compressedFileInfo, err := videoCompress.CompressFileToSize(files[0], "aboba", int64(48 * 1024 * 1024))
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

		_, err = telegram.SendVideo(of, caption[0], channelName[0], compressedFileInfo.CompressedFilename)
		//defer compressedFileInfo.Reader.Close()
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
	} else {
		fmt.Println("small file")
		fmt.Println(files[0].Filename)
		of, err := files[0].Open()
		if err != nil {
			c.JSON(400, gin.H{"error":err})
			c.Abort()
			return
		}
		reader := io.Reader(of)
		fmt.Println(caption[0])
		fmt.Println(channelName[0])
		_, err = telegram.SendVideoBytes(reader, files[0].Filename, caption[0], channelName[0])
		if err != nil {
			c.JSON(400, gin.H{"error":err})
			c.Abort()
			return
		}
		of.Close()
	}

	c.JSON(200, gin.H{"message": "success"})
}

func sendVideoNoteHandler(c *gin.Context) {
	multipart, _ := c.MultipartForm()
	files := multipart.File["videoNote"]
	channelName := multipart.Value["channelName"]
	for _, file := range files {
		telegram.SendVideoNote(file, channelName[0])
	}
	c.JSON(200, gin.H{"message": "success"})
}

//func sendVideoNoteHandler(c *gin.Context) {
//	multipart, _ := c.MultipartForm()
//	files := multipart.File["videoNote"]
//	channelName := multipart.Value["channelName"]
//
//	for _, file := range files {
//		videoData, err := resizeVideoToSquare(file)
//		if err != nil {
//			c.JSON(500, gin.H{"error": "Failed to resize video"})
//			return
//		}
//
//		telegram.SendVideoNote(videoData, channelName[0])
//	}
//
//	c.JSON(200, gin.H{"message": "success"})
//}

//func resizeVideoToSquare(file *multipart.FileHeader) ([]byte, error) {
//	src, err := file.Open()
//	if err != nil {
//		fmt.Println("error 1")
//		return nil, err
//	}
//	defer src.Close()
//
//	// Read the original video data
//	videoData, err := io.ReadAll(src)
//	if err != nil {
//		fmt.Println("error 2")
//		return nil, err
//	}
//
//	cmd := exec.Command("ffmpeg", "-i", "-", "-vf", "scale=384:384", "-c:v", "libx264", "-f", "mp4", "pipe:1")
//	cmd.Stdin = bytes.NewReader(videoData)
//
//	var stderr bytes.Buffer
//	cmd.Stderr = &stderr
//
//	resizedVideoData, err := cmd.Output()
//	if err != nil {
//		fmt.Println("Error executing FFmpeg:", err)
//		fmt.Println("FFmpeg error output:", stderr.String())
//		return nil, err
//	}
//
//	return resizedVideoData, nil
//}

func sendLocationHandler(c *gin.Context) {
	var body struct {
		Latitude  string `json:"latitude"`
		Longitude string `json:"longitude"`
		ChatId    string `json:"channelName"`
	}
	err := c.Bind(&body)
	if err != nil {
		fmt.Println(err)
		return
	}

	telegram.SendLocation(body.Latitude, body.Longitude, body.ChatId)

	c.JSON(200, gin.H{
		"statusCode": "success",
	})
}

func sendVenueHandler(c *gin.Context) {
	var body struct {
		Latitude  string `json:"latitude"`
		Longitude string `json:"longitude"`
		Title     string `json:"title"`
		Address   string `json:"address"`
		ChatId    string `json:"channelName"`
	}
	err := c.Bind(&body)
	if err != nil {
		fmt.Println(err)
		return
	}

	telegram.SendVenue(body.Latitude, body.Longitude, body.Title, body.Address, body.ChatId)

	c.JSON(200, gin.H{
		"statusCode": "success",
	})
}

//func postDelete(c *gin.Context) {
//
//	id := c.Param("id")
//
//	utils.DB.Delete(&models.Post{}, id)
//
//	c.JSON(200, gin.H{
//		"statusCode": "success",
//	})
//}
func sendPhotoHandler(c *gin.Context) {

	multipart, _ := c.MultipartForm()
	files := multipart.File["photo"]
	caption := multipart.Value["caption"]
	for _, file := range files {
		of, err := file.Open()
		if err != nil {
			c.JSON(400, gin.H{"error":err})
			c.Abort()
			return
		}
		defer of.Close()
		telegram.SendPhoto(of, caption[0], file.Filename)
	}
	c.JSON(200, gin.H{"aboba":"aboba"})
}
func sendMediaGroupHandler(c *gin.Context) {
	multipart, _ := c.MultipartForm()
	files := multipart.File["media"]
	caption := multipart.Value["caption"]
	var filenames []string
	var fileList []*io.Reader
	for _, file := range files {
		of, err := file.Open()
		if err != nil {
			c.JSON(400, gin.H{"error":err})
			c.Abort()
			return
		}
		defer of.Close()
		readerPtr := io.Reader(of)
		fileList = append(fileList, &readerPtr)
		filenames = append(filenames, file.Filename)
	}
	//_, err := telegram.SendMediaGroupLazy(files, caption[0])
	_, err := telegram.SendMediaGroup(fileList, filenames, caption[0])
	if err!= nil {
		c.JSON(500, gin.H{"error":err})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{"message":"success"})
}

func sendMediaGroupLinks(c *gin.Context) {
	_, _ = telegram.SendMediaGroupLinks([]string{"https://i.pinimg.com/originals/61/3b/f8/613bf893ab736ac25c6f6dde1bbacc4a.jpg", "https://sweetpeaskitchen.com/wp-content/uploads/2020/06/Choc-Rasp-Cheesecake-7-1-scaled.jpg"}, "text")
	c.JSON(200, gin.H{"message":"success"})
}
func sendDiceHandler(c *gin.Context)  {
	telegram.SendDice()
	c.JSON(200, gin.H{"a":"a"})
}
