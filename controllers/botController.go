package controllers

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golearn/api/telegram"
	"golearn/models"
	mongoRepository "golearn/repository"
	"golearn/utils/auth"
	"golearn/utils/jsonHelper"
	"io"
	"os"
)

func SetupBotRoutes(r *gin.Engine) {
	botGroup := r.Group("/bot")
	//botGroup.GET("/getMe", getMeHandler)

	botGroup.Use(auth.SubLevelMiddleware(1))

	botGroup.POST("/sendMessage", jsonHelper.MakeHttpHandler(sendMessageHandler))
	botGroup.POST("/sendPhoto", jsonHelper.MakeHttpHandler(sendPhotoHandler))
	botGroup.POST("/sendMediaGroup", jsonHelper.MakeHttpHandler(sendMediaGroupHandler))
	botGroup.POST("/sendAudio", jsonHelper.MakeHttpHandler(sendAudioHandler))
	botGroup.POST("/sendVoice", jsonHelper.MakeHttpHandler(sendVoiceHandler))
	botGroup.POST("/sendVideo", jsonHelper.MakeHttpHandler(sendVideoHandler))
	botGroup.POST("/sendVideoNote", jsonHelper.MakeHttpHandler(sendVideoNoteHandler))
	botGroup.POST("/sendLocation", jsonHelper.MakeHttpHandler(sendLocationHandler))
	botGroup.POST("/sendVenue", jsonHelper.MakeHttpHandler(sendVenueHandler))
}

//
//func getMeHandler(c *gin.Context) {
//	telegram.GetMe()
//	c.JSON(200, gin.H{"message": "success"})
//}

type SendMessageRequest struct {
	Text        string `json:"text"`
	ChannelName string `json:"channelName"`
	BotToken string `json:"botToken"`
}

// @Summary Send message
// @Tags bot
// @Description Send text message to some channel
// @ID SendMessage
// @Accept json
// @Produce json
// @Param request body controllers.SendMessageRequest true "Message body"
// @Success 200 {string} a
// @Failure 404 {object} jsonHelper.ApiError "Channel with given ID does not exist"
// @Failure 400 {object} jsonHelper.ApiError "Error identifying user"
// @Failure 417 {object} jsonHelper.ApiError "Error identifying user"
// @Failure 500 {object} jsonHelper.ApiError
// @Failure default {object} jsonHelper.ApiError
// @Router /bot/sendMessage [post]
func sendMessageHandler(c *gin.Context) error {
	var body SendMessageRequest
	err := c.Bind(&body)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}

	channelRepo := mongoRepository.NewChannelRepo()
	chID, err := uuid.Parse(body.ChannelName)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    "Internal server error",
			Status: 500,
		}
	}

	channel, err:=channelRepo.FindByID(context.Background(), chID)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    "Channel with specified ID does not exist",
			Status: 404,
		}
	}

	err = telegram.SendMessage(body.BotToken,body.Text, channel.Name)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}

	c.JSON(200, gin.H{"message": "Message has been scheduled successfully"})
	return nil
}

// @Summary Send audio
// @Tags bot
// @Description Send text message to some channel
// @ID SendAudio
// @Accept mpfd
// @Produce json
// @Param caption formData string true "Text of post"
// @Param botToken formData string true "Token of bot that manages given channels"
// @Param channelName formData string true "Channel ID"
// @Param audio formData file true "Audio message file"
// @Success 200 {string} string "OK"
// @Failure 500 {object} jsonHelper.ApiError
// @Failure default {object} jsonHelper.ApiError
// @Router /bot/sendAudio [post]
func sendAudioHandler(c *gin.Context) error {
	multipart, _ := c.MultipartForm()
	files := multipart.File["audio"]
	caption := multipart.Value["caption"]
	channelName := multipart.Value["channelName"]
	botToken := multipart.Value["botToken"]

	channelRepo := mongoRepository.NewChannelRepo()
	chID, err := uuid.Parse(channelName[0])
	if err != nil {
		return jsonHelper.ApiError{
			Err:    "Internal server error",
			Status: 500,
		}
	}

	channel, err:=channelRepo.FindByID(context.Background(), chID)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    "Channel with specified ID does not exist",
			Status: 404,
		}
	}


	for _, file := range files {
		_, err := telegram.SendAudio(botToken[0],file, caption[0], channel.Name)
		if err != nil {
			return jsonHelper.ApiError{
				Err:    err.Error(),
				Status: 500,
			}
		}
	}
	c.JSON(200, gin.H{"message": "success"})
	return nil
}

// @Summary Send voice
// @Tags bot
// @Description Send voice message to some channel
// @ID SendVoice
// @Accept mpfd
// @Produce json
// @Param caption formData string true "Text of post"
// @Param voice formData file true "Voice message file"
// @Param channelName formData string true "Channel ID"
// @Param botToken formData string true "Token of bot that manages given channels"
// @Success 200 {string} string "OK"
// @Failure 500 {object} jsonHelper.ApiError
// @Failure default {object} jsonHelper.ApiError
// @Router /bot/sendVoice [post]
func sendVoiceHandler(c *gin.Context) error {
	multipart, _ := c.MultipartForm()
	files := multipart.File["voice"]
	caption := multipart.Value["caption"]
	channelName := multipart.Value["channelName"]
	botToken := multipart.Value["botToken"]

	channelRepo := mongoRepository.NewChannelRepo()
	chID, err := uuid.Parse(channelName[0])
	if err != nil {
		return jsonHelper.ApiError{
			Err:    "Internal server error",
			Status: 500,
		}
	}

	channel, err:=channelRepo.FindByID(context.Background(), chID)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    "Channel with specified ID does not exist",
			Status: 404,
		}
	}

	for _, file := range files {
		_, err := telegram.SendVoice(botToken[0],file, caption[0], channel.Name)
		if err != nil {
			return jsonHelper.ApiError{
				Err:    err.Error(),
				Status: 500,
			}
		}
	}
	c.JSON(200, gin.H{"message": "success"})
	return nil
}

// @Summary Send video
// @Tags bot
// @Description Send video message to some channel
// @ID SendVideo
// @Accept mpfd
// @Produce json
// @Param caption formData string true "Text of post"
// @Param video formData file true "Voice message file"
// @Param channelName formData string true "Channel ID"
// @Param botToken formData string true "Token of bot that manages given channels"
// @Success 200 {string} string "OK"
// @Failure 400 {object} jsonHelper.ApiError
// @Failure 500 {object} jsonHelper.ApiError
// @Failure default {object} jsonHelper.ApiError
// @Router /bot/sendVideo [post]
func sendVideoHandler(c *gin.Context) error {
	multipart, _ := c.MultipartForm()
	files := multipart.File["video"]
	caption := multipart.Value["caption"]
	channelName := multipart.Value["channelName"]
	botToken := multipart.Value["botToken"]

	channelRepo := mongoRepository.NewChannelRepo()
	chID, err := uuid.Parse(channelName[0])
	if err != nil {
		return jsonHelper.ApiError{
			Err:    "Internal server error",
			Status: 500,
		}
	}

	channel, err:=channelRepo.FindByID(context.Background(), chID)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    "Channel with specified ID does not exist",
			Status: 404,
		}
	}
	of, err := files[0].Open()
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 400,
		}
	}
	reader := io.Reader(of)
	_, err = telegram.SendVideoBytes(botToken[0],reader, files[0].Filename, caption[0], channel.Name)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 400,
		}
	}
	defer of.Close()

	c.JSON(200, gin.H{"message": "success"})
	return nil
}


// @Summary Send video
// @Tags bot
// @Description Send video message to some channel
// @ID SendVideoNote
// @Accept mpfd
// @Produce json
// @Param videoNote formData file true "Voice message file"
// @Param channelName formData string true "Channel ID"
// @Success 200 {string} string "OK"
// @Failure 400 {object} jsonHelper.ApiError
// @Failure 500 {object} jsonHelper.ApiError
// @Failure default {object} jsonHelper.ApiError
// @Router /bot/sendVideoNote [post]
func sendVideoNoteHandler(c *gin.Context) error {
	multipart, _ := c.MultipartForm()
	files := multipart.File["videoNote"]
	channelName := multipart.Value["channelName"]
	botToken := multipart.Value["botToken"]

	channelRepo := mongoRepository.NewChannelRepo()
	chID, err := uuid.Parse(channelName[0])
	if err != nil {
		return jsonHelper.ApiError{
			Err:    "Internal server error",
			Status: 500,
		}
	}

	channel, err:=channelRepo.FindByID(context.Background(), chID)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    "Channel with specified ID does not exist",
			Status: 404,
		}
	}

	for _, file := range files {
		_, err := telegram.SendVideoNote(botToken[0],file, channel.Name)
		if err != nil {
			return jsonHelper.ApiError{
				Err:    err.Error(),
				Status: 500,
			}
		}
	}
	c.JSON(200, gin.H{"message": "success"})
	return nil
}


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

type SendLocationRequest struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	ChatId    string `json:"channelName"`
	BotToken string `json:"botToken"`
}

// @Summary Send location
// @Tags bot
// @Description Send location message to some channel
// @ID SendLocation
// @Accept mpfd
// @Produce json
// @Param request body controllers.SendLocationRequest true "Location body"
// @Success 200 {string} string "OK"
// @Failure 400 {object} jsonHelper.ApiError
// @Failure 500 {object} jsonHelper.ApiError
// @Failure default {object} jsonHelper.ApiError
// @Router /bot/sendLocation [post]
func sendLocationHandler(c *gin.Context) error {
	var body SendLocationRequest
	err := c.Bind(&body)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 400,
		}
	}

	channelRepo := mongoRepository.NewChannelRepo()
	chID, err := uuid.Parse(body.ChatId)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    "Internal server error",
			Status: 500,
		}
	}

	channel, err:=channelRepo.FindByID(context.Background(), chID)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    "Channel with specified ID does not exist",
			Status: 404,
		}
	}



	err = telegram.SendLocation(body.BotToken,body.Latitude, body.Longitude, channel.Name)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 400,
		}
	}

	c.JSON(200, gin.H{
		"statusCode": "success",
	})
	return nil
}

type SendVenueRequest struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Title     string `json:"title"`
	Address   string `json:"address"`
	ChatId    string `json:"channelName"`
	BotToken string `json:"botToken"`
}

// @Summary Send venue
// @Tags bot
// @Description Send venue
// @ID SendVenue
// @Accept mpfd
// @Produce json
// @Param request body controllers.SendVenueRequest true "Venue request body"`
// @Success 200 {string} string "OK"
// @Failure 400 {object} jsonHelper.ApiError
// @Failure 500 {object} jsonHelper.ApiError
// @Failure default {object} jsonHelper.ApiError
// @Router /bot/sendVenue [post]
func sendVenueHandler(c *gin.Context) error {
	var body SendVenueRequest
	err := c.Bind(&body)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 400,
		}
	}

	channelRepo := mongoRepository.NewChannelRepo()
	chID, err := uuid.Parse(body.ChatId)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    "Internal server error",
			Status: 500,
		}
	}

	channel, err:=channelRepo.FindByID(context.Background(), chID)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    "Channel with specified ID does not exist",
			Status: 404,
		}
	}


	err = telegram.SendVenue(body.BotToken,body.Latitude, body.Longitude, body.Title, body.Address, channel.Name)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}

	c.JSON(200, gin.H{
		"status": "success",
	})
	return nil
}


// @Summary Send photo
// @Tags bot
// @Description Send photo message to some channel
// @ID SendPhoto
// @Accept mpfd
// @Produce json
// @Param caption formData string true "Text of post"
// @Param photo formData file true "Photo message file"
// @Param channelName formData string true "Channel ID"
// @Success 200 {string} string "OK"
// @Failure 400 {object} jsonHelper.ApiError
// @Failure 404 {object} jsonHelper.ApiError "Channel with specified ID does not exist"
// @Failure 500 {object} jsonHelper.ApiError
// @Failure default {object} jsonHelper.ApiError
// @Router /bot/sendPhoto [post]
func sendPhotoHandler(c *gin.Context) error {
	multipart, _ := c.MultipartForm()
	files := multipart.File["photo"]
	caption := multipart.Value["caption"]
	channelName := multipart.Value["channelName"]
	botToken := multipart.Value["botToken"]

	channelRepo := mongoRepository.NewChannelRepo()
	chID, err := uuid.Parse(channelName[0])
	if err != nil {
		return jsonHelper.ApiError{
			Err:    "Internal server error",
			Status: 500,
		}
	}

	channel, err:=channelRepo.FindByID(context.Background(), chID)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    "Channel with specified ID does not exist",
			Status: 404,
		}
	}

	for _, file := range files {
		of, err := file.Open()
		if err != nil {
			return jsonHelper.ApiError{
				Err:    err.Error(),
				Status: 400,
			}
		}
		defer of.Close()
		_, err = telegram.SendPhoto(botToken[0],of, caption[0], file.Filename, channel.Name)
		if err != nil {
			return jsonHelper.ApiError{
				Err:    err.Error(),
				Status: 500,
			}
		}
	}
	c.JSON(200, gin.H{"status": "success"})
	return nil
}

// @Summary Send media
// @Tags bot
// @Description Send mediagroup message to some channel
// @ID SendMediaGroup
// @Accept mpfd
// @Produce json
// @Param caption formData string true "Text of post"
// @Param media formData file true "Media message file (max of 10 files are allowed)"
// @Param chat formData string true "ID of channel"
// @Param fileTypes formData string true "Stringified map of each file type (Example: \"[{\"fileName1\":\"video\"},{\"fileName2\":\"photo\"}]\")"
// @Success 200 {string} a
// @Failure 400 {object} jsonHelper.ApiError
// @Failure 500 {object} jsonHelper.ApiError
// @Failure default {object} jsonHelper.ApiError
// @Router /bot/sendPhoto [post]
func sendMediaGroupHandler(c *gin.Context) error {
	multipart, _ := c.MultipartForm()
	files := multipart.File["media"]
	caption := multipart.Value["caption"]
	chat := multipart.Value["chat"]
	botToken := multipart.Value["botToken"]

	fileTypesField := multipart.Value["fileTypes"]

	var data []map[string]string

	if err := json.Unmarshal([]byte(fileTypesField[0]), &data); err != nil {
		return jsonHelper.ApiError{
			Err:    "Error processing file types",
			Status: 400,
		}
	}

	fileTypeMap := make(map[string]string)

	for _, entry := range data {
		for filename, fileType := range entry {
			fileTypeMap[filename] = fileType
		}
	}

	channelRepo := mongoRepository.NewChannelRepo()
	chID, err := uuid.Parse(chat[0])
	if err != nil {
		return jsonHelper.ApiError{
			Err:    "Internal server error",
			Status: 500,
		}
	}

	channel, err:=channelRepo.FindByID(context.Background(), chID)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    "Channel with specified ID does not exist",
			Status: 404,
		}
	}


	var filenames []string
	var fileList []*io.Reader
	var fileModels []models.File

	for _, file := range files {
		of, err := file.Open()
		defer of.Close()
		if err != nil {
			return jsonHelper.ApiError{
				Err:    err.Error(),
				Status: 500,
			}
		}
		fileID, err := uuid.NewRandom()
		fileModel := models.File{
			ID:         fileID,
			BucketName: os.Getenv("mediaGroupBucketName"),
			Type:       fileTypeMap[file.Filename],
			PostID:     uuid.UUID{},
		}
		readerPtr := io.Reader(of)
		fileList = append(fileList, &readerPtr)
		filenames = append(filenames, file.Filename)
		fileModels = append(fileModels, fileModel)
	}
	//_, err := telegram.SendMediaGroupLazy(files, caption[0])
	_, err = telegram.SendMediaGroup(botToken[0],fileList, filenames,fileModels, caption[0], channel.Name)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}
	c.JSON(200, gin.H{"message": "success"})
	return nil
}