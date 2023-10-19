package controllers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golearn/models"
	mongoRepository "golearn/repository"
	"golearn/utils/auth"
	"golearn/utils/jsonHelper"
	"golearn/utils/s3"
	"golearn/utils/videoCompress"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func SetupScheduleRoutes(r *gin.Engine)  {
	scheduleGroup := r.Group("/schedule")

	scheduleGroup.Use(auth.AuthMiddleware)

	//scheduleGroup.Use(auth.SubscriptionMiddleware())
	//scheduleGroup.Use(auth.SubLevelMiddleware(0))
	scheduleGroup.POST("/message", jsonHelper.MakeHttpHandler(scheduleMessageHandler))
	scheduleGroup.POST("/photo", jsonHelper.MakeHttpHandler(schedulePhotoHandler))
	scheduleGroup.POST("/mediaGroup", scheduleMediaGroupHandler)
	scheduleGroup.POST("/video", scheduleVideoHandler)
	scheduleGroup.POST("/audio", scheduleAudioHandler)
	scheduleGroup.POST("/voice", scheduleVoiceHandler)
	scheduleGroup.GET("/", getScheduledPostHandler)
	scheduleGroup.GET("/:id", jsonHelper.MakeHttpHandler(getPostHandler))
	scheduleGroup.DELETE("/delete/:id", jsonHelper.MakeHttpHandler(deletePostHandler))
	scheduleGroup.GET("/image/:imageName", jsonHelper.MakeHttpHandler(getPostImage) )
	scheduleGroup.GET("/date/:scheduled", jsonHelper.MakeHttpHandler(getPostsByDate))
}

func getPostsByDate(c *gin.Context) error {

	postsRepo := mongoRepository.NewPostRepository()

	authUserEmail, exists := c.Get("userEmail")
	if !exists {
		return jsonHelper.ApiError{
			Err:    "User unauthorized",
			Status: 401,
		}
	}

	scheduledTime, err := time.Parse( "2006-01-15",c.Param("scheduled"))

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
		case posts:=<-postsch:
			c.JSON(200, gin.H{"posts":posts})
			c.Abort()
			return nil
		}

	}
}

func getPostImage(c *gin.Context) error {

	postsRepo := mongoRepository.NewPostRepository()

	authUserEmail, exists := c.Get("userEmail")
	if !exists {
		return jsonHelper.ApiError{
			Err:    "User unauthorized",
			Status: 401,
		}
	}
	user, err := mongoRepository.GetUserByEmail(fmt.Sprintf("%s", authUserEmail))
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}
	imageName, err := uuid.Parse(c.Param("imageName"))

	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 400,
		}
	}
	post ,err := postsRepo.GetPostByImageName(context.Background(), imageName)
	if user.ID.String() != post.ID.String() {

	}
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
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

func deletePostHandler(c *gin.Context) error {
	postsRepo := mongoRepository.NewPostRepository()

	authUserEmail, exists := c.Get("userEmail")
	if !exists {
		return jsonHelper.ApiError{
			Err:    "User unauthorized",
			Status: 401,
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
	c.JSON(200, gin.H{"status":deleted})
	return nil
}

func getPostHandler(c *gin.Context) error {

	postsRepo := mongoRepository.NewPostRepository()

	authUserEmail, exists := c.Get("userEmail")
	if !exists {
		return jsonHelper.ApiError{
			Err:    "User unauthorized",
			Status: 401,
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
	c.JSON(200, gin.H{"post":post})
	return nil
}

func getScheduledPostHandler(c *gin.Context) {
	postsRepo := mongoRepository.NewPostRepository()
	authUserEmail, exists := c.Get("userEmail")
	if !exists {
		c.JSON(401, gin.H{"Error":"Unauthorized"})
		c.Abort()
		return
	}
	postsch := make(chan []models.Post, 1)
	user, err := mongoRepository.GetUserByEmail(fmt.Sprintf("%s", authUserEmail))
	if err != nil {
		fmt.Println("aboba error")
		c.JSON(500, gin.H{"Error":err})
		c.Abort()
		return
	}
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	go postsRepo.GetPostsByUserID(ctx, user.ID, wg, postsch)
	wg.Wait()
	for {
		select {
			case posts:= <-postsch:
				c.JSON(200, gin.H{"posts":posts})
				c.Abort()
				return
			case <-ctx.Done():
				fmt.Println("aboba error 1")
				c.JSON(500, gin.H{"Error":"Timeout. Due to internal server error"})
				c.Abort()
				return
		}
	}
}

func scheduleVoiceHandler(c *gin.Context)  {
	multipart, _ := c.MultipartForm()
	files := multipart.File["audio"]
	caption := multipart.Value["caption"]
	title := multipart.Value["title"]
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
		Title: title[0],
		ID: postID,
		Text:        caption[0],
		ChannelName: channelName[0],
		Type:        "voice",
		Files:       []uuid.UUID{postID},
		Scheduled:   parsedTime,
	}
	err = mongoRepository.SaveScheduledPost(&post)
	if err != nil {
		c.JSON(400, gin.H{"error":err})
		c.Abort()
		return
	}
	if err != nil {
		c.JSON(400, gin.H{"error":err})
		c.Abort()
		return
	}
	err = s3.LoadAudio(postID.String(), &file)
	if err != nil {
		c.JSON(400, gin.H{"error":err.Error()})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{"message": "success"})
}

func scheduleAudioHandler(c *gin.Context) {
	multipart, _ := c.MultipartForm()
	files := multipart.File["audio"]
	title := multipart.Value["title"]
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
		Title: title[0],
		Text:        caption[0],
		ChannelName: channelName[0],
		Type:        "audio",
		Files:       []uuid.UUID{postID},
		Scheduled:   parsedTime,
	}
	err = mongoRepository.SaveScheduledPost(&post)
	if err != nil {
		c.JSON(400, gin.H{"error":err})
		c.Abort()
		return
	}
	err = s3.LoadAudio(postID.String(), &file)
	if err != nil {
		c.JSON(400, gin.H{"error":err.Error()})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{"message": "success"})
}

func scheduleMessageHandler(c *gin.Context) error {

	postRepo := mongoRepository.NewPostRepository()

	var body struct {
		Title string `json:"title"`
		Text string `json:"text"`
		Chat string `json:"chat"`
		Time string `json:"scheduled"`
	}
	jsonHelper.BindWithException(&body, c)
	postId, _ := uuid.NewRandom()

	authUserEmail, exists := c.Get("userEmail")
	if !exists {
		return jsonHelper.ApiError{
			Err:    "User unauthorized",
			Status: 401,
		}
	}
	user, err := mongoRepository.GetUserByEmail(fmt.Sprintf("%s", authUserEmail))
	parsedTime, _ := time.Parse("2006 01-02 15:04 -0700 MST", body.Time)
	post:=models.Post{
		Title: body.Title,
		Text:        body.Text,
		ChannelName: body.Chat,
		Type:        "message",
		UserID:      user.ID,
		Files:    	 nil,
		ID: 		 postId,
		Scheduled: parsedTime,
	}
	fmt.Println(postId)
	err = postRepo.SavePostWithId(&post, postId)
	if err != nil {
		fmt.Println(err)
		return jsonHelper.ApiError{
			Err:    "Error saving post",
			Status: 500,
		}
	}
	c.JSON(200, gin.H{"post":post})
	return nil
}

func schedulePhotoHandler(c *gin.Context) error {
	multipart, _ := c.MultipartForm()

	files := multipart.File["photo"]
	title := multipart.Value["title"]
	caption := multipart.Value["caption"]
	channelName := multipart.Value["channelName"]
	scheduledTime := multipart.Value["scheduled"]
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
			Err:    err.Error(),
			Status: 400,
		}
	}
	authUserEmail, exists := c.Get("userEmail")
	if !exists {
		return jsonHelper.ApiError{
			Err:    "User unauthorized",
			Status: 401,
		}
	}
	user, err := mongoRepository.GetUserByEmail(fmt.Sprintf("%s", authUserEmail))
	post := models.Post{
		Title: title[0],
		ID: postID,
		UserID: user.ID,
		Text:        caption[0],
		ChannelName: channelName[0],
		Type:        "photo",
		Files:       []uuid.UUID{postID},
		Scheduled:   parsedTime,
	}
	err = mongoRepository.SaveScheduledPost(&post)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    fmt.Sprintf("Error saving post : %s", err.Error()),
			Status: 500,
		}
	}
	err = s3.LoadImage(context.Background(), postID.String(), &file)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    fmt.Sprintf("Error loading image : %s", err.Error()),
			Status: 500,
		}
	}
	c.JSON(200, gin.H{"message":"success"})
	return nil
}

func scheduleMediaGroupHandler(c *gin.Context)  {
	multipart, _ := c.MultipartForm()
	files := multipart.File["media"]
	title := multipart.Value["title"]
	caption := multipart.Value["caption"]
	channelName := multipart.Value["chat"]
	scheduledTime := multipart.Value["scheduled"]
	parsedTime, _ := time.Parse("2006 01-02 15:04 -0700 MST", scheduledTime[0])
	fmt.Println(parsedTime)
	postID, err := uuid.NewRandom()
	if err != nil {
		c.JSON(400, gin.H{"error":err})
		c.Abort()
		return
	}
	post := models.Post{
		Title: title[0],
		ID: postID,
		Text:        caption[0],
		ChannelName: channelName[0],
		Type:        "mediaGroup",
		Files:       []uuid.UUID{postID},
		Scheduled:   parsedTime,
	}
	err = mongoRepository.SaveScheduledPost(&post)
	if err != nil {
		c.JSON(400, gin.H{"error":err})
		c.Abort()
		return
	}
	var fileList []uuid.UUID
	for _, file := range files{
		fmt.Println("processing file")
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
			fmt.Println(err)
			c.JSON(400, gin.H{"error":err})
			c.Abort()
			return
		}
		fmt.Println(fileID)
		fileList = append(fileList, fileID)
		err = of.Close()
		if err != nil {
			c.JSON(400, gin.H{"error":err})
			c.Abort()
			return
		}
	}
	err = mongoRepository.UpdateFilesList(postID, fileList)
	if err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{"error":err})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{})
}

func scheduleVideoHandler(c *gin.Context)  {
	multipart, _ := c.MultipartForm()
	files := multipart.File["video"]
	title := multipart.Value["title"]
	caption := multipart.Value["caption"]
	channelName := multipart.Value["channelName"]
	scheduledTime := multipart.Value["scheduled"]

	parsedTime, _ := time.Parse("2006 01-02 15:04 -0700 MST", scheduledTime[0])
	postID, err := uuid.NewRandom()
	if err != nil {
		c.JSON(400, gin.H{"error":err})
		c.Abort()
		return
	}
	post := models.Post{
		ID: postID,
		Title: title[0],
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