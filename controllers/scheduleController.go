package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"golearn/models"
	mongoRepository "golearn/repository"
	"golearn/utils/auth"
	"golearn/utils/jsonHelper"
	"golearn/utils/s3"
	"os"
	"sync"
	"time"
)

func SetupScheduleRoutes(r *gin.Engine) {
	scheduleGroup := r.Group("/schedule")

	scheduleGroup.GET("/image/:imageName", jsonHelper.MakeHttpHandler(getPostImage))
	scheduleGroup.GET("/video/:videoName", jsonHelper.MakeHttpHandler(getPostsVideo))

	scheduleGroup.Use(auth.AuthMiddleware)

	//scheduleGroup.Use(auth.SubscriptionMiddleware())
	//scheduleGroup.Use(auth.SubLevelMiddleware(0))
	scheduleGroup.POST("/message", jsonHelper.MakeHttpHandler(scheduleMessageHandler))
	scheduleGroup.POST("/photo", jsonHelper.MakeHttpHandler(schedulePhotoHandler))
	scheduleGroup.POST("/mediaGroup", jsonHelper.MakeHttpHandler(scheduleMediaGroupHandler))
	scheduleGroup.POST("/video", jsonHelper.MakeHttpHandler(scheduleVideoHandler))
	scheduleGroup.POST("/audio", jsonHelper.MakeHttpHandler(scheduleAudioHandler))
	scheduleGroup.POST("/voice", jsonHelper.MakeHttpHandler(scheduleVoiceHandler))
	scheduleGroup.GET("/", jsonHelper.MakeHttpHandler(getScheduledPostHandler))
	scheduleGroup.GET("/:id", jsonHelper.MakeHttpHandler(getPostHandler))
	scheduleGroup.DELETE("/delete/:id", jsonHelper.MakeHttpHandler(deletePostHandler))
	scheduleGroup.GET("/date/:scheduled", jsonHelper.MakeHttpHandler(getPostsByDate))
}

// @Summary Get posts videos
// @Tags posts
// @Description Receive post's video banner
// @ID GetPostsVideo
// @Accept json
// @Produce octet-stream
// @Param videoName path string true "Name of the image"
// @Success 200 {string} a
// @Failure 500 {object} jsonHelper.ApiError
// @Failure default {object} jsonHelper.ApiError
// @Router /schedule/video/{videoName} [get]
func getPostsVideo(c *gin.Context) error {
	videoName, err := uuid.Parse(c.Param("videoName"))

	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 400,
		}
	}

	image, err := s3.GetVideo(videoName.String())
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}

	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}
	c.DataFromReader(200, -1, "application/octet-stream", image, nil)
	return nil
}

// @Summary Get posts by date
// @Tags posts
// @Description Receive posts by date
// @ID GetPostsByDate
// @Accept json
// @Produce json
// @Param scheduled query string true "Date"
// @Success 200 {string} a
// @Failure 500 {object} jsonHelper.ApiError
// @Failure default {object} jsonHelper.ApiError
// @Router /schedule/date/{scheduled} [get]
func getPostsByDate(c *gin.Context) error {

	postsRepo := mongoRepository.NewPostRepository()

	authUserEmail, exists := c.Get("userEmail")
	if !exists {
		return jsonHelper.ApiError{
			Err:    "User unauthorized",
			Status: 417,
		}
	}

	scheduledTime, err := time.Parse("2006-01-15", c.Param("scheduled"))

	if err != nil {
		return jsonHelper.ApiError{
			Err:    "Error parsing this huynia",
			Status: 400,
		}
	}

	user, err := mongoRepository.GetUserByEmail(fmt.Sprintf("%s", authUserEmail))
	if err != nil {
		return jsonHelper.ApiError{
			Err:    "User does not exist",
			Status: 404,
		}
	}

	postsch := make(chan []models.Post, 1)
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	go postsRepo.GetPostsByDate(ctx, scheduledTime, user.ID, wg, postsch)
	wg.Wait()

	for {
		select {
		case <-ctx.Done():
			return jsonHelper.ApiError{
				Err:    "Timeout",
				Status: 500,
			}
		case posts := <-postsch:
			c.JSON(200, gin.H{"posts": posts})
			c.Abort()
			return nil
		}

	}
}

// @Summary Get post's image
// @Tags posts
// @Description Receive post's image banner
// @ID GetPostsImage
// @Accept json
// @Produce json
// @Param imageName path string true "Name of the image"
// @Success 200 {string} a
// @Failure 500 {object} jsonHelper.ApiError
// @Failure default {object} jsonHelper.ApiError
// @Router /schedule/image/{imageName} [post]
func getPostImage(c *gin.Context) error {
	imageName, err := uuid.Parse(c.Param("imageName"))

	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 400,
		}
	}

	image, err := s3.GetImage(imageName.String())
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}

	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}
	//c.Stream(200, "application/octet-stream",func(w io.Writer) (int,error) {
	//
	//})
	c.DataFromReader(200, -1, "application/octet-stream", image, nil)
	return nil
}

// @Summary Delete post by id
// @Tags posts
// @Description Delete post by id
// @ID DeletePost
// @Accept json
// @Produce json
// @Param id path string true "ID of post to delete"
// @Success 200 {string} a
// @Failure 400, 417 {object} jsonHelper.ApiError
// @Failure 500 {object} jsonHelper.ApiError
// @Failure default {object} jsonHelper.ApiError
// @Router /schedule/delete/{id} [delete]
func deletePostHandler(c *gin.Context) error {
	postsRepo := mongoRepository.NewPostRepository()

	authUserEmail, exists := c.Get("userEmail")
	if !exists {
		return jsonHelper.ApiError{
			Err:    "User unauthorized",
			Status: 417,
		}
	}
	user, err := mongoRepository.GetUserByEmail(fmt.Sprintf("%s", authUserEmail))
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}
	postID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}

	post, err := postsRepo.GetPostByID(context.Background(), postID)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}
	if post.UserID.String() != user.ID.String() {
		return jsonHelper.ApiError{
			Err:    "",
			Status: 403,
		}
	}
	deleted := postsRepo.DeletePostByID(context.Background(), post.ID)
	c.JSON(200, gin.H{"status": deleted})
	return nil
}

// @Summary Receive post by id
// @Tags posts
// @Description Receive post's image banner
// @ID GetPost
// @Accept json
// @Produce json
// @Param id path string true "ID of post to receive"
// @Success 200 {string} a
// @Failure 400, 417 {object} jsonHelper.ApiError
// @Failure 500 {object} jsonHelper.ApiError
// @Failure default {object} jsonHelper.ApiError
// @Router /schedule/{id} [get]
func getPostHandler(c *gin.Context) error {

	postsRepo := mongoRepository.NewPostRepository()

	authUserEmail, exists := c.Get("userEmail")
	if !exists {
		return jsonHelper.ApiError{
			Err:    "User unauthorized",
			Status: 417,
		}
	}
	user, err := mongoRepository.GetUserByEmail(fmt.Sprintf("%s", authUserEmail))
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}
	postID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}
	post, err := postsRepo.GetPostByID(context.Background(), postID)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}
	if post.UserID.String() != user.ID.String() {
		return jsonHelper.ApiError{
			Err:    "",
			Status: 403,
		}
	}
	c.JSON(200, gin.H{"post": post})
	return nil
}

// @Summary Get all scheduled posts
// @Tags posts
// @Description Receive all scheduled posts
// @ID GetScheduledPosts
// @Accept json
// @Produce json
// @Success 200 {string} a
// @Failure 400, 417 {object} jsonHelper.ApiError
// @Failure 500 {object} jsonHelper.ApiError
// @Failure default {object} jsonHelper.ApiError
// @Router /schedule/ [get]
func getScheduledPostHandler(c *gin.Context) error {
	postsRepo := mongoRepository.NewPostRepository()
	authUserEmail, exists := c.Get("userEmail")
	if !exists {
		return jsonHelper.ApiError{
			Err:    "Unauthorized",
			Status: 417,
		}
	}
	postsch := make(chan []models.Post, 1)
	user, err := mongoRepository.GetUserByEmail(fmt.Sprintf("%s", authUserEmail))
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	go postsRepo.GetPostsByUserID(ctx, user.ID, wg, postsch)
	wg.Wait()
	for {
		select {
		case posts := <-postsch:
			c.JSON(200, gin.H{"posts": posts})
			c.Abort()
			return nil
		case <-ctx.Done():
			return jsonHelper.ApiError{
				Err:    "Timeout. Due to internal server error",
				Status: 500,
			}
		}
	}
}

// @Summary Schedule voice
// @Tags posts
// @Description Receive all scheduled posts
// @ID ScheduleVoice
// @Accept json
// @Produce json
// @Param caption body string true "Text of post"
// @Param audio body file true "Voice message file"
// @Param channelName body string true "Channel name"
// @Param title body string true "Title of post (in-app only, won't affect telegram)"
// @Param scheduled body string true "Scheduled date"
// @Success 200 {string} a
// @Failure 400, 417 {object} jsonHelper.ApiError
// @Failure 500 {object} jsonHelper.ApiError
// @Failure default {object} jsonHelper.ApiError
// @Router /schedule/voice [post]
func scheduleVoiceHandler(c *gin.Context) error {
	multipart, _ := c.MultipartForm()
	files := multipart.File["audio"]
	caption := multipart.Value["caption"]
	title := multipart.Value["title"]
	channelName := multipart.Value["channelName"]
	scheduledTime := multipart.Value["scheduled"]

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

	parsedTime, _ := time.Parse("2006 01-02 15:04 -0700 UTC", scheduledTime[0])
	file, err := files[0].Open()
	defer file.Close()
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 400,
		}
	}
	postID, err := uuid.NewRandom()
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 400,
		}
	}
	post := models.Post{
		Title:       title[0],
		ID:          postID,
		Text:        caption[0],
		ChannelID: channel.ID,
		Type:        "voice",
		Files:       []uuid.UUID{postID},
		Scheduled:   parsedTime,
	}
	err = mongoRepository.SaveScheduledPost(context.Background(),&post)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 400,
		}
	}
	err = s3.LoadAudio(postID.String(), &file)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}

	c.JSON(200, gin.H{"message": "success"})
	return nil
}

// @Summary Schedule audio
// @Tags posts
// @Description Schedule audio
// @ID ScheduleAudio
// @Accept json
// @Produce json
// @Param caption body string true "Text of post"
// @Param audio body file true "Audio message file"
// @Param channelName body string true "Channel name"
// @Param title body string true "Title of post (in-app only, won't affect telegram)"
// @Param scheduled body string true "Scheduled date"
// @Success 200 {string} a
// @Failure 400, 417 {object} jsonHelper.ApiError
// @Failure 500 {object} jsonHelper.ApiError
// @Failure default {object} jsonHelper.ApiError
// @Router /schedule/audio [post]
func scheduleAudioHandler(c *gin.Context) error {

	fileRepo := mongoRepository.NewFileRepo()

	multipart, _ := c.MultipartForm()
	files := multipart.File["audio"]
	title := multipart.Value["title"]
	caption := multipart.Value["caption"]
	channelName := multipart.Value["channelName"]
	scheduledTime := multipart.Value["scheduled"]

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

	parsedTime, _ := time.Parse("2006 01-02 15:04 -0700 UTC", scheduledTime[0])
	file, err := files[0].Open()
	defer file.Close()
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 400,
		}
	}
	postID, err := uuid.NewRandom()
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}
	fileID, err := uuid.NewRandom()
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}
	post := models.Post{
		ID:          postID,
		Title:       title[0],
		Text:        caption[0],
		ChannelID: channel.ID,
		Type:        "audio",
		Files:       []uuid.UUID{postID},
		Scheduled:   parsedTime,
	}
	savedFile := models.File{
		ID:         fileID,
		BucketName: os.Getenv("audioBucketName"),
		Type:       "audio",
		PostID:     postID,
	}

	err = fileRepo.Save(context.Background(), &savedFile)

	if err != nil {
		return jsonHelper.ApiError{
			Err:    "Error saving file",
			Status: 500,
		}
	}

	err = mongoRepository.SaveScheduledPost(context.Background(),&post)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}
	err = s3.LoadAudio(postID.String(), &file)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}

	c.JSON(200, gin.H{"message": "success"})
	return nil
}

type ScheduleMessageRequest struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	Chat  string `json:"chat"`
	Time  string `json:"scheduled"`
}

// @Summary Schedule message
// @Tags posts
// @Description Schedule message
// @ID ScheduleMessage
// @Accept json
// @Produce json
// @Param request body controllers.ScheduleMessageRequest true "Request body"
// @Success 200 {string} a
// @Failure 400, 417 {object} jsonHelper.ApiError
// @Failure 500 {object} jsonHelper.ApiError
// @Failure default {object} jsonHelper.ApiError
// @Router /schedule/audio [post]
func scheduleMessageHandler(c *gin.Context) error {

	postRepo := mongoRepository.NewPostRepository()

	var body ScheduleMessageRequest
	jsonHelper.BindWithException(&body, c)
	postId, _ := uuid.NewRandom()

	authUserEmail, exists := c.Get("userEmail")
	if !exists {
		return jsonHelper.ApiError{
			Err:    "User unauthorized",
			Status: 417,
		}
	}
	user, err := mongoRepository.GetUserByEmail(fmt.Sprintf("%s", authUserEmail))
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 417,
		}
	}

	channelRepo := mongoRepository.NewChannelRepo()
	chID, err := uuid.Parse(body.Chat)
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

	parsedTime, _ := time.Parse("2006 01-02 15:04 -0700 MST", body.Time)
	post := models.Post{
		Title:       body.Title,
		Text:        body.Text,
		ChannelID: channel.ID,
		Type:        "message",
		UserID:      user.ID,
		Files:       nil,
		ID:          postId,
		Scheduled:   parsedTime,
	}
	err = postRepo.SavePostWithId(&post, postId)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    "Error saving post",
			Status: 500,
		}
	}
	c.JSON(200, gin.H{"post": post})
	return nil
}

// @Summary Schedule photo
// @Tags posts
// @Description Schedule photo
// @ID SchedulePhoto
// @Accept json
// @Produce json
// @Param caption body string true "Text of post"
// @Param photo body file true "photo message file"
// @Param channelName body string true "Channel name"
// @Param title body string true "Title of post (in-app only, won't affect telegram)"
// @Param scheduled body string true "Scheduled date"
// @Success 200 {string} a
// @Failure 400, 417 {object} jsonHelper.ApiError
// @Failure 500 {object} jsonHelper.ApiError
// @Failure default {object} jsonHelper.ApiError
// @Router /schedule/photo [post]
func schedulePhotoHandler(c *gin.Context) error {
	multipart, _ := c.MultipartForm()
	fileRepo:=mongoRepository.NewFileRepo()

	files := multipart.File["photo"]
	title := multipart.Value["title"]
	caption := multipart.Value["caption"]
	channelName := multipart.Value["channelName"]
	scheduledTime := multipart.Value["scheduled"]


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

	parsedTime, _ := time.Parse("2006 01-02 15:04 -0700 MST", scheduledTime[0])
	file, err := files[0].Open()
	defer file.Close()
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 400,
		}
	}
	postID, err := uuid.NewRandom()
	if err != nil {
		return jsonHelper.ApiError{
			Err:    "Internal server error",
			Status: 500,
		}
	}
	fileID, err := uuid.NewRandom()
	if err != nil {
		return jsonHelper.ApiError{
			Err:    "Internal server error",
			Status: 500,
		}
	}
	authUserEmail, exists := c.Get("userEmail")
	if !exists {
		return jsonHelper.ApiError{
			Err:    "User unauthorized",
			Status: 417,
		}
	}
	user, err := mongoRepository.GetUserByEmail(fmt.Sprintf("%s", authUserEmail))
	post := models.Post{
		Title:       title[0],
		ID:          postID,
		UserID:      user.ID,
		Text:        caption[0],
		ChannelID:   channel.ID,
		Type:        "photo",
		Files:       []uuid.UUID{postID},
		Scheduled:   parsedTime,
	}
	savedFile := models.File{
		ID:         fileID,
		BucketName: os.Getenv("imageBucketName"),
		Type:       "photo",
		PostID:     postID,
	}

	err = mongoRepository.WithTransaction(context.Background(), func(ctx mongo.SessionContext) (interface{}, error) {
		err = mongoRepository.SaveScheduledPost(ctx,&post)
		if err != nil {
			return nil,err
		}

		err = fileRepo.Save(ctx, &savedFile)

		err = s3.LoadImage(context.Background(), postID.String(), &file)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	c.JSON(200, gin.H{"message": "success"})
	return nil
}

// @Summary Schedule mediagroup
// @Tags posts
// @Description Schedule mediagroup
// @ID ScheduleMediaGroup
// @Accept json
// @Produce json
// @Param caption body string true "Text of post"
// @Param media body file true "Media message file"
// @Param channelName body string true "Channel name"
// @Param title body string true "Title of post (in-app only, won't affect telegram)"
// @Param scheduled body string true "Channel name"
// @Success 200 {string} a
// @Failure 400, 417 {object} jsonHelper.ApiError
// @Failure 500 {object} jsonHelper.ApiError
// @Failure default {object} jsonHelper.ApiError
// @Router /schedule/mediaGroup [post]
func scheduleMediaGroupHandler(c *gin.Context) error {

	fileRepo := mongoRepository.NewFileRepo()

	authUserEmail, exists := c.Get("userEmail")
	if !exists {
		return jsonHelper.ApiError{
			Err:    "User unauthorized",
			Status: 417,
		}
	}
	user, err := mongoRepository.GetUserByEmail(fmt.Sprintf("%s", authUserEmail))
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 417,
		}
	}

	// Processing request
	multipart, _ := c.MultipartForm()
	files := multipart.File["media"]
	title := multipart.Value["title"]
	caption := multipart.Value["caption"]
	channelName := multipart.Value["chat"]
	scheduledTime := multipart.Value["scheduled"]

	fileTypesField := multipart.Value["fileTypes"]

	var data []map[string]string

	if err := json.Unmarshal([]byte(fileTypesField[0]), &data); err != nil {
		return jsonHelper.ApiError{
			Err:    "Error processing file types",
			Status: 0,
		}
	}

	fileTypeMap := make(map[string]string)

	for _, entry := range data {
		for filename, fileType := range entry {
			fileTypeMap[filename] = fileType
		}
	}

	// done processing request

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

	parsedTime, _ := time.Parse("2006 01-02 15:04 -0700 MST", scheduledTime[0])
	postID, err := uuid.NewRandom()
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 400,
		}
	}
	post := models.Post{
		Title:       title[0],
		ID:          postID,
		Text:        caption[0],
		ChannelID:   channel.ID,
		Type:        "mediaGroup",
		Files:       []uuid.UUID{postID},
		Scheduled:   parsedTime,
		UserID:      user.ID,
	}
	var fileList []uuid.UUID
	err = mongoRepository.WithTransaction(context.Background(), func(ctx mongo.SessionContext) (interface{}, error) {
		err = mongoRepository.SaveScheduledPost(context.Background(),&post)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			of, err := file.Open()
			if err != nil {
				return nil, err
			}

			fileID, err := uuid.NewRandom()
			if err != nil {
				return nil,err
			}

			savedFile := models.File{
				ID:         fileID,
				BucketName: os.Getenv("mediaGroupBucketName"),
				Type:       fileTypeMap[file.Filename],
				PostID: post.ID,
			}
			err = fileRepo.Save(context.Background(),&savedFile)


			err = s3.LoadMedia(context.Background(), fileID.String(), &of)
			if err != nil {
				return nil, err
			}
			fileList = append(fileList, fileID)
			err = of.Close()
			if err != nil {
				return nil, err
			}
		}
		return nil, nil
	})
	err = mongoRepository.UpdateFilesList(context.Background(), postID, fileList)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}
	c.JSON(200, gin.H{})
	return nil
}

// @Summary Schedule video
// @Tags posts
// @Description Schedule video
// @ID ScheduleVideo
// @Accept json
// @Produce json
// @Param caption body string true "Text of post"
// @Param video body file true "Video message file"
// @Param channelName body string true "Channel name"
// @Param title body string true "Title of post (in-app only, won't affect telegram)"
// @Param scheduled body string true "Scheduled date"
// @Success 200 {string} a
// @Failure 400, 417 {object} jsonHelper.ApiError
// @Failure 500 {object} jsonHelper.ApiError
// @Failure default {object} jsonHelper.ApiError
// @Router /schedule/video [post]
func scheduleVideoHandler(c *gin.Context) error {

	fileRepo := mongoRepository.NewFileRepo()

	authUserEmail, exists := c.Get("userEmail")
	if !exists {
		return jsonHelper.ApiError{
			Err:    "User unauthorized",
			Status: 417,
		}
	}
	user, err := mongoRepository.GetUserByEmail(fmt.Sprintf("%s", authUserEmail))
	multipart, _ := c.MultipartForm()
	files := multipart.File["video"]
	title := multipart.Value["title"]
	caption := multipart.Value["caption"]
	channelName := multipart.Value["channelName"]
	scheduledTime := multipart.Value["scheduled"]

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

	parsedTime, _ := time.Parse("2006 01-02 15:04 -0700 MST", scheduledTime[0])
	postID, err := uuid.NewRandom()
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}
	fileID ,err := uuid.NewRandom()
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}
	post := models.Post{
		ID:          postID,
		Title:       title[0],
		Text:        caption[0],
		ChannelID: channel.ID,
		Type:        "video",
		Files:       []uuid.UUID{},
		Scheduled:   parsedTime,
		UserID: user.ID,
	}
	file := models.File{
		ID: fileID,
		BucketName: os.Getenv("videoBucketName"),
		Type:       "video",
		PostID:     postID,
	}
	err = fileRepo.Save(context.Background(), &file)
	err = mongoRepository.SaveScheduledPost(context.Background(),&post)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 400,
		}
	}
	var fileIDList []uuid.UUID

	of, err := files[0].Open()
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 400,
		}
	}
	defer of.Close()
	err = s3.LoadVideoMultipart(fileID.String(), &of)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 400,
		}
	}
	fileIDList = append(fileIDList, fileID)

	err = mongoRepository.UpdateFilesList(context.Background(),postID, fileIDList)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 400,
		}
	}
	c.JSON(200, gin.H{})
	return nil
}