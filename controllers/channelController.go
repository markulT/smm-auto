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
)

func SetupChannelRoutes(r *gin.Engine) {
	channelGroup := r.Group("/channel")
	channelGroup.Use(auth.AuthMiddleware)

	channelGroup.DELETE("/delete/:id", jsonHelper.MakeHttpHandler(deleteChannelHandler))
	channelGroup.GET("/", jsonHelper.MakeHttpHandler(getAllChannels))

	channelGroup.POST("/add", jsonHelper.MakeHttpHandler(addChannelHandler))
	channelGroup.POST("/assignToken", jsonHelper.MakeHttpHandler(assignTokenToChannelHandler))
}

type AddChannelRequest struct {
	ChannelName string `json:"channelName"`
}

func getAllChannels(c *gin.Context) error {

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

	channelRepo:=mongoRepository.NewChannelRepo()
	channelList, err := channelRepo.GetAllByUserID(context.Background(), user.ID)

	if err != nil {
		return jsonHelper.ApiError{
			Err:    "Interval server error while fetching channels",
			Status: 500,
		}
	}

	c.JSON(200, gin.H{"channels":channelList})
	return nil
}

func addChannelHandler(c *gin.Context) error {

	var body AddChannelRequest
	jsonHelper.BindWithException(&body, c)

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
	chID, err := uuid.NewRandom()
	fmt.Println(chID)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    "Internal server error",
			Status: 500,
		}
	}
	channel := &models.Channel{
		Name:             body.ChannelName,
		UserID: user.ID,
	}

	err = channelRepo.SaveNewChannel(context.Background(), channel)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    "Internal server error on saving your channel",
			Status: 500,
		}
	}
	c.JSON(200, gin.H{"channel":channel})
	return nil
}

func deleteChannelHandler(c *gin.Context) error {

	channelIDParam := c.Param("id")
	channelID, err := uuid.Parse(channelIDParam)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    "Wrong channel id passed (can't parse)",
			Status: 400,
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
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 417,
		}
	}

	channelRepo := mongoRepository.NewChannelRepo()
	channel, err := channelRepo.FindByID(context.Background(), channelID)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    "No channel with such ID",
			Status: 404,
		}
	}

	if channel.UserID.String() != user.ID.String() {
		return jsonHelper.ApiError{
			Err:    "You don't have access to this channel",
			Status: 403,
		}
	}
	
	err = channelRepo.DeleteChannel(context.Background(), channelID)

	if err != nil {
		return jsonHelper.ApiError{
			Err:    "Internal server error on deleting channel",
			Status: 500,
		}
	}
	c.JSON(200, gin.H{})
	return nil
}

type AssignTokenToChannelRequest struct {
	Token string `json:"token"`
	ChannelID string `json:"channelId"`
}

func assignTokenToChannelHandler(c *gin.Context) error {

	var body AssignTokenToChannelRequest
	jsonHelper.BindWithException(&body, c)

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
	chID , err := uuid.Parse(body.ChannelID)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    "Can't parse id",
			Status: 500,
		}
	}
	channel, err := channelRepo.FindByID(context.Background(), chID)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    "Channel does not exist",
			Status: 404,
		}
	}

	if user.ID.String() != channel.UserID.String() {
		if err != nil {
			return jsonHelper.ApiError{
				Err:    "Channel does not exist",
				Status: 404,
			}
		}
	}

	err = channelRepo.AssignBotToken(context.Background(), body.Token, chID)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    "Error assigning token to bot",
			Status: 500,
		}
	}

	return nil
}

