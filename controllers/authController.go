package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"golearn/models"
	"golearn/utils"
	"golearn/utils/auth"
	"golearn/utils/jsonHelper"
	"os"
)

func SetupAuthRoutes(r *gin.Engine) {
	authGroup := r.Group("/auth")
	authGroup.POST("/login", login)
	authGroup.POST("/refresh", refresh)
	authGroup.POST("/signup", signup)
}

type LoginRequestBody struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

func login(c *gin.Context) {
	fmt.Println("Login request is here")
	var body struct{
		Email string `json:"email"`
		Password string `json:"password"`
	}

	jsonHelper.BindWithException(&body, c)

	var userExists models.User
	err := utils.DB.Where("email = ?", body.Email).First(&userExists).Error
	if err!=nil {
		c.JSON(404, gin.H{"error":"user does not exist"})
		return
	}
	var userFromDB models.User
	utils.DB.Where("email = ?", body.Email).First(&userFromDB)
	if err := bcrypt.CompareHashAndPassword([]byte(userFromDB.Password), []byte(body.Password));err!=nil {
		c.JSON(403, gin.H{"error":"Wrong pass"})
		c.Abort()
		return
	}
	tokens := auth.GenerateTokens(map[string]interface{}{
		"email":userFromDB.Email,
		"password": userFromDB.Password,
	}, c)

	c.JSON(200, gin.H{
		"accessToken": tokens.AccessToken,
		"refreshToken": tokens.RefreshToken,
		"email": userFromDB.Email,
	})

}


func refresh(c *gin.Context) {
	secretKey := os.Getenv("secretKeyRefresh")
	var body struct {
		RefreshToken string `json:"refreshToken"`
	}
	jsonHelper.BindWithException(&body, c)
	var userFromDb models.User
	email, err := auth.GetSubject(body.RefreshToken, secretKey)
	if err!=nil {
		fmt.Println(err)
	}
	fmt.Println("Email is : ")
	fmt.Println(email)
	utils.DB.Where("email = ?", email).First(&userFromDb)
	if _, err := auth.Validate(body.RefreshToken, secretKey);err!=nil {
		c.JSON(400, gin.H{"message":err})
		c.Abort()
		return
	}
	tokens := auth.GenerateTokens(map[string]interface{}{
		"email":userFromDb.Email,
	}, c)
	c.JSON(200, gin.H{
		"accessToken":tokens.AccessToken,
		"refreshToken":tokens.RefreshToken,
		"email": userFromDb.Email,
	})
}



func signup(c *gin.Context) {
	var body struct{
		Email string `json:"email"`
		Password string `json:"password"`
	}
	jsonHelper.BindWithException(&body, c)
	auth.UserExists(body.Email, c)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err!=nil {
		return
	}
	newUser := models.User{Email: body.Email, Password: string(hashedPassword)}
	tokens := auth.GenerateTokens(map[string]interface{}{
		"email":    newUser.Email,
	}, c)

	if err := utils.DB.Create(&newUser).Error; err!=nil {
		c.JSON(400, gin.H{"message":"Error creating user"})
		return
	}
	c.JSON(200, gin.H{
		"accessToken": tokens.AccessToken,
		"refreshToken": tokens.RefreshToken,
		"email":newUser.Email,
	})
}