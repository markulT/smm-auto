package auth

import (
	"github.com/gin-gonic/gin"
	"os"
	"strings"
)

func AuthMiddleware(c *gin.Context) gin.HandlerFunc {
	return func(context *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(401, gin.H{"error":"Unauthorized"})
			c.Abort()
			return
		}
		accessToken := authHeader[7:]
		if _, err := Validate(accessToken, os.Getenv("secretKey"));err!=nil {
			c.JSON(401, gin.H{"error":"Invalid token"})
		}
		c.Next()
	}
}
