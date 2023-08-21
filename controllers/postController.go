package controllers

import (
	"github.com/gin-gonic/gin"
	"golearn/models"
	"golearn/utils"
	"reflect"
)

func PostCreate(c *gin.Context) {

	var body struct{
		body string
		title string
	}
	c.Bind(&body)

	post := models.Post{Title: body.title, Body: body.body}

	result := utils.DB.Create(&post)

	if result.Error != nil {
		c.Status(408)
		return
	}

	c.JSON(200, gin.H{
		"post":post,
	})
}

func PostIndex(c *gin.Context)  {
	var posts []models.Post
	utils.DB.Find(&posts)
	c.JSON(200, gin.H{
		"posts": posts,
	})
}

func PostShow(c *gin.Context) {
	id:=c.Param("id")
	var post models.Post
	utils.DB.First(&post, id)
	c.JSON(200, gin.H{
		"post":post,
	})
}

func PostUpdateWithReflect(c *gin.Context) {

	var body struct{
		id string
		Data map[string]interface{}
	}
	if err:=c.Bind(&body);err!=nil {
		c.JSON(400, gin.H{"error":"Invalid JSON"})
		return
	}
	var post models.Post
	if err:=utils.DB.First(&post,body.id).Error;err!=nil {
		c.JSON(400, gin.H{"error":"Invalid JSON at the request"})
		return
	}
	for fieldName, fieldValue := range body.Data{
		field := reflect.ValueOf(&post).Elem().FieldByName(fieldName)
		if field.IsValid() && field.CanSet() {
			field.Set(reflect.ValueOf(fieldValue))
		}
	}
	utils.DB.Save(&post)
	c.JSON(200, gin.H{
		"post":post,
	})
}

func PostDelete(c *gin.Context) {
	id:=c.Param("id")
	var post models.Post
	utils.DB.First(&post, id)
	utils.DB.Delete(&models.Post{}, id)
	c.JSON(200, gin.H{
		"post":post,
		"status":"deleted",
	})
}
