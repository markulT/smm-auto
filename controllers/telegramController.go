package controllers
//
//import (
//	"fmt"
//	"github.com/gin-gonic/gin"
//)
//
//func SetupTelegramRoutes(r *gin.Engine)  {
//	telegramGroup := r.Group("/telegram")
//	telegramGroup.POST("/sendMessage", sendMessageHandler)
//	telegramGroup.POST("/confirmCode", confirmCodeHandler)
//}
//
//
//
//func confirmCodeHandler(c *gin.Context) {
//	var body struct {
//		Phone string `json:"phone"`
//		Code string `json:"code"`
//		PhoneCodeHash string `json:"phoneCodeHash"`
//	}
//	err := c.Bind(&body)
//	fmt.Println(body.Code)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//}