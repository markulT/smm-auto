package auth

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"golearn/repository"
	"net/http"
)

func UserExists(email string, c *gin.Context) {
	//var count int64
	//utils.DB.Model(&models.User{}).Where("email = ?", email).Count(&count)
	_, err := repository.GetUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		c.Abort()
		return
	}
}
func GenerateTokens(body map[string]interface{}, c *gin.Context) (Tokens) {
	var tokens Tokens
	accessToken, accessErr := CreateAccessToken(body)
	if accessErr != nil {
		c.JSON(400, gin.H{"message":"Error generating access token"})
	}
	refreshToken, refreshErr := CreateRefreshToken(body)
	if refreshErr != nil {
		c.JSON(400, gin.H{"message":"Error generating refresh token"})
	}
	tokens.AccessToken = accessToken
	tokens.RefreshToken = refreshToken
	return tokens
}
func ComparePasswords(plainPassword string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}
