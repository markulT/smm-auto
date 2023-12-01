package controllers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"golearn/models"
	mongoRepository "golearn/repository"
	"golearn/utils/auth"
	"golearn/utils/jsonHelper"
)

func SetupArchiveRoutes(r *gin.Engine) {
	archiveGroup := r.Group("/archive")

	archiveGroup.Use(auth.AuthMiddleware)

	archiveGroup.GET("/", jsonHelper.MakeHttpHandler(getAllArchivedPosts))

}

type ArchivedPosts struct {
	Posts []ResponsePost `json:"posts"`
}

type ResponsePost struct {
	models.Post
	Files []models.File `bson:"files" json:"files"`
}

// @Summary Get all archived posts
// @Tags archive
// @Description Get all archived posts
// @ID GetAllArchivedPosts
// @Accept json
// @Produce json
// @Success 200 {object} ArchivedPosts
// @Failure 417 {object} jsonHelper.ApiError "Error identifying user"
// @Failure 500 {object} jsonHelper.ApiError "Internal server error"
// @Failure default {object} jsonHelper.ApiError
// @Router /archive/ [get]
func getAllArchivedPosts(c *gin.Context) error {
	var posts []models.Post
	postsRepo := mongoRepository.NewPostRepository()
	fileRepo := mongoRepository.NewFileRepo()
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
	var responsePostList []*ResponsePost
	for _, post := range posts {
		files, err := fileRepo.FindManyByIDList(context.Background(), post.Files)
		if err != nil {
			return jsonHelper.ApiError{
				Err:    "Failed to load filetypes",
				Status: 500,
			}
		}
		rp := ResponsePost{post, files}
		responsePostList = append(responsePostList, &rp)
	}
	c.JSON(200, gin.H{"posts":responsePostList})
	return nil
}
