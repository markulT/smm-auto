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

	channelGroup.POST("/add", jsonHelper.MakeHttpHandler(addChannelHandler))
	channelGroup.POST("/assignToken", jsonHelper.MakeHttpHandler(assignTokenToChannelHandler))
	channelGroup.GET("/", jsonHelper.MakeHttpHandler(getAllChannelsHandler))
	//channelGroup.GET("/:id", jsonHelper.MakeHttpHandler(getAllChannelsHandler))
}

type AddChannelRequest struct {
	ChannelName string `json:"channelName"`
}

type GetAllChannelsResponse struct {
	Channels []models.Channel `json:"channels"`
}

// @Summary Get all channels
// @Tags channel
// @Description Returns an array of Channel objects
// @ID GetAllChannels
// @Accept json
// @Produce json
// @Success 200 {object} GetAllChannelsResponse
// @Failure 400 {object} jsonHelper.ApiError "Error identifying user"
// @Failure 417 {object} jsonHelper.ApiError "Error identifying user"
// @Failure 500 {object} jsonHelper.ApiError "Internal server error"
// @Failure default {object} jsonHelper.ApiError
// @Router /channel/ [get]
func getAllChannelsHandler(c *gin.Context) error {

	//var channelsList []models.Channel

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
			Err:    "Error identifying user",
			Status: 417,
		}
	}

	channelRepo := mongoRepository.NewChannelRepo()
	channelsList, err := channelRepo.FindAllByUserID(context.Background(), user.ID)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    "Internal server error. Error loading user's channels",
			Status: 500,
		}
	}

	c.JSON(200, gin.H{"channels": *channelsList})
	return nil
}


// @Summary Add channel
// @Tags channel
// @Description Adds channel to user's channelList
// @ID AddChannel
// @Accept json
// @Produce json
// @Success 200 {string} a
// @Failure 400 {object} jsonHelper.ApiError "Error identifying user"
// @Failure 417 {object} jsonHelper.ApiError "Error identifying user"
// @Failure 500 {object} jsonHelper.ApiError
// @Failure default {object} jsonHelper.ApiError
// @Router /channel/add [post]
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
	if err != nil {
		return jsonHelper.ApiError{
			Err:    "Internal server error",
			Status: 500,
		}
	}
	chID, err := uuid.NewUUID()
	if err != nil {
		return jsonHelper.ApiError{
			Err:    "Error generating channel's ID, try again",
			Status: 500,
		}
	}
	channel := &models.Channel{
		Name:             body.ChannelName,
		UserID: user.ID,
		ID: chID,
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

// @Summary Delete channel
// @Tags channel
// @Description Remove channel from user's list
// @ID DeleteChannel
// @Accept json
// @Produce json
// @Param id path string true "channel ID"
// @Success 200 {string} string "OK"
// @Failure 400 {object} jsonHelper.ApiError "Error identifying user"
// @Failure 417 {object} jsonHelper.ApiError "Error identifying user"
// @Failure 403 {object} jsonHelper.ApiError "User does not have access to this channel"
// @Failure 500 {object} jsonHelper.ApiError "Internal server error"
// @Failure default {object} jsonHelper.ApiError
// @Router /channel/add [post]
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
	c.JSON(200, gin.H{})
	return nil
}

