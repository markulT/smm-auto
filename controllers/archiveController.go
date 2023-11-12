package controllers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"golearn/models"
	mongoRepository "golearn/repository"
	"golearn/utils/jsonHelper"
)

func SetupArchiveRoutes(r *gin.Engine) {
	archiveGroup := r.Group("/archive")

	archiveGroup.GET("/", jsonHelper.MakeHttpHandler(GetAllArchivedPosts))

}

func GetAllArchivedPosts(c *gin.Context) error {
	var posts []models.Post
	postsRepo := mongoRepository.NewPostRepository()
	authUserEmail, exists := c.Get("userEmail")
	if !exists {
		return jsonHelper.ApiError{
			Err:    "Unauthorized",
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
	posts, err = postsRepo.GetAllArchivedPostsByUserID(context.Background(), user.ID)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}
	c.JSON(200, gin.H{"posts":posts})
	return nil
}
