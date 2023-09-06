package jsonHelper

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func BindWithException(body interface{}, c *gin.Context) {
	if err:= c.Bind(body); err!=nil {
		fmt.Println(err)
		c.JSON(400, gin.H{"error":"invalid JSON"})
		c.Abort()
		return
	}
}
