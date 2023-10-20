package auth

import (
	"github.com/gin-gonic/gin"
	"strings"
)

func AuthMiddleware(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(401, gin.H{"error":"Unauthorized"})
			c.Abort()
			return
		}
		accessToken := authHeader[7:]
		if _, err := Validate(accessToken);err!=nil {
			c.JSON(401, gin.H{"error":"Invalid token"})
		}
		userEmail, err := GetSubject(accessToken)
		if err != nil {
			c.JSON(401, gin.H{"error":"Error extracting subject from token (invalid token)"})
		}
		c.Set("userEmail", userEmail)
		c.Next()
}
