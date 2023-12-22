package controllers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	mongoRepository "golearn/repository"
	"golearn/utils/analytics"
	"golearn/utils/auth"
	"golearn/utils/jsonHelper"
	"strconv"
	"time"
)

type analyticsController struct {
	analyticsService analytics.AnalyticsService
	analyticsRepo mongoRepository.AnalyticsRepo
}

func SetupAnalyticsRoutes(r *gin.Engine, as analytics.AnalyticsService, analyticsRepo mongoRepository.AnalyticsRepo) {
	ac := analyticsController{analyticsService: as}
	ac.analyticsService = as
	ac.analyticsRepo = analyticsRepo
	analyticsGroup := r.Group("/analytics")
	analyticsGroup.Use(auth.AuthMiddleware)

	analyticsGroup.GET("/", jsonHelper.MakeHttpHandler(ac.getChannelsViews))

}

func (ac *analyticsController) getChannelsViews(c *gin.Context) error {

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

	timestampStart := c.Query("timestampStart")

	unixTimeStart, err := strconv.ParseInt(timestampStart, 10, 64)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    "Wrong date parameter in request",
			Status: 400,
		}
	}
	dateStart := time.Unix(unixTimeStart, 0)

	timestampEnd := c.Query("timestampEnd")

	unixTimeEnd, err := strconv.ParseInt(timestampEnd, 10, 64)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    "Wrong date parameter in request",
			Status: 400,
		}
	}

	// Convert the Unix timestamp to a time.Time object
	dateEnd := time.Unix(unixTimeEnd, 0)
	channelName := c.Query("channelName")
	fmt.Println(channelName)
	fmt.Println("cholera")
	analyticsItems, err := ac.analyticsRepo.GetManyByChannelIDAndDate(context.Background(), channelName, dateStart, dateEnd )
	if err != nil {
		return jsonHelper.ApiError{
			Err:    "Error occured while searching for analytics items",
			Status: 500,
		}
	}

	fmt.Println(user)

	c.JSON(200, gin.H{"analyticsItems":analyticsItems})
	return nil
}



