package utils

import (
	"golearn/api/telegram"
	"golearn/models"
	"time"
)

func FetchAndProcessPosts() {
	for range time.Tick(time.Second * 10) {
		go func() {
			var posts []models.Post
			if err := DB.Find(&posts).Error; err != nil {
				return
			}
			for post := 0; post < len(posts); post++ {
				loc, _ := time.LoadLocation(posts[post].TimeZone)

				now := time.Now().In(loc)

				scheduledTime, err := time.ParseInLocation("2006-01-02 15:04", posts[post].Scheduled, loc)

				if err != nil {
					return
				}

				if scheduledTime.Sub(now) <= 0 && posts[post].Status == "scheduled" {
					var postToSent models.Post
					DB.First(&postToSent, posts[post].ID)

					DB.Model(&postToSent).Updates(models.Post{
						Status: "sent",
					})
					telegram.SendMessage(posts[post].Text, posts[post].ChannelName)
				}
			}
		}()
	}
}
