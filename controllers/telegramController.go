package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golearn/api/telegram/auth"
	"golearn/api/telegram/messages"
)

func SetupTelegramRoutes(r *gin.Engine)  {
	telegramGroup := r.Group("/telegram")
	telegramGroup.POST("/sendMessage", sendMessageHandler)
	telegramGroup.POST("/sendCode", sendCodeHandler)
	telegramGroup.POST("/confirmCode", confirmCodeHandler)
}

func sendMessageHandler(c *gin.Context) {
	var body struct {
		Text string `json:"text"`
	}
	err := c.Bind(&body)
	fmt.Println(body.Text)
	if err != nil {
		fmt.Println(err)
		return
	}
	messages.SendMessage(body.Text)
	c.JSON(200, gin.H{"status":"Sent"})
}
func sendCodeHandler(c *gin.Context) {
	var body struct {
		Phone string `json:"phone"`
	}
	err := c.Bind(&body)
	fmt.Println(body.Phone)
	if err != nil {
		fmt.Println(err)
		return
	}
	result := auth.SendCode(body.Phone)
	c.JSON(200, gin.H{
		"data": result,
	})
}

func confirmCodeHandler(c *gin.Context) {
	var body struct {
		Phone string `json:"phone"`
		Code string `json:"code"`
		PhoneCodeHash string `json:"phoneCodeHash"`
	}
	err := c.Bind(&body)
	fmt.Println(body.Code)
	if err != nil {
		fmt.Println(err)
		return
	}
	auth.ConfirmCode(body.Code, body.Phone, body.PhoneCodeHash)
}