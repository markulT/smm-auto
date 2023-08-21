package controllers

import (
	"github.com/gin-gonic/gin"
	"golearn/models"
	"golearn/utils"
	"golearn/utils/auth"
	"golearn/utils/jsonHelper"
)

func SetupAuthRoutes(r *gin.Engine) {
	authGroup := r.Group("/auth")
	authGroup.POST("/login", login)
	authGroup.POST("/refresh", refresh)
	authGroup.POST("/signup", signup)
}

func login(c *gin.Context) {
	var body struct{
		email string
		password string
	}
	jsonHelper.BindWithException(body, c)
	var userExists models.User
	err := utils.DB.Where("email = ?", body.email).First(&userExists).Error
	if err!=nil {
		c.JSON(400, gin.H{"error":"invalid JSON"})
		return
	}


}
func refresh(c *gin.Context) {

}
func signup(c *gin.Context) {
	var body struct{
		email string
		password string
	}
	jsonHelper.BindWithException(body, c)

	if userExists := utils.DB.Where("email = ?", body.email).First(&models.User{}); userExists!=nil {
		c.JSON(400, gin.H{"error":"User already exists"})
		return
	}

	newUser := models.User{Email: body.email, Password: body.password}
	token, err := auth.CreateAccessToken(map[string]interface{}{
		"email":    newUser.Email,
		"password": newUser.Password,
	})
	if err != nil {
		return
	}
	utils.DB.Save(newUser)
	c.JSON(200, gin.H{
		"accessToken": token,
	})
}