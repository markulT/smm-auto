package scheduler

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"golearn/api/telegram"
	"golearn/models"
	"golearn/repository"
	"golearn/utils"
	"golearn/utils/notifications"
	"golearn/utils/s3"
	"io"
	"os"
	"strconv"
	"sync"
	"time"
)

func FetchAndProcessPosts() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			ProcessPosts()
		}
	}
}

func ProcessPosts() {
	fmt.Print("processingPosts")
	numElements := int(getTableLength("posts"))
	batchSize, _ := strconv.Atoi(os.Getenv("batchSize"))
	var wg sync.WaitGroup
	for start := 0; start < numElements; start += batchSize {
		end := start + batchSize - 1
		if end >= numElements {
			end = numElements - 1
		}
		wg.Add(1)
		go processBatch(start, end, &wg)
	}
	wg.Wait()
}

func getTableLength(tableName string) int64 {
	collection := utils.DB.Collection(tableName)

	count, err := collection.CountDocuments(context.Background(), bson.D{})
	if err != nil {
		return 0
	}
	return count
}

func processBatch(start, end int, wg *sync.WaitGroup) {
	defer wg.Done()

	var scheduledPosts []models.Post

	notificationService := notifications.DefaultNotificationService{}

	scheduledPosts = *repository.GetScheduledPostRelations(start, end-start+1, false)

	for _, scheduledPost := range scheduledPosts {
		originalTimezone, offset := scheduledPost.Scheduled.Zone()
		currentTime := time.Now().In(time.FixedZone(originalTimezone, offset))
		if scheduledPost.Scheduled.Before(currentTime) {
			switch scheduledPost.Type {
			case "message":
				telegram.SendMessage(scheduledPost.Text, scheduledPost.ChannelName)
				err := repository.ArchivizePost(scheduledPost.ID)
				if err != nil {
					return
				}
				notificationService.SendNotification("Notification", "Scheduled message sent!", scheduledPost.DeviceToken)
			case "photo":
				image, err := s3.GetImage(scheduledPost.Files[0].String())
				if err != nil {
					return
				}
				_, err = telegram.SendPhoto(image, scheduledPost.Text, scheduledPost.Files[0].String(), scheduledPost.ChannelName)
				if err != nil {
					return
				}
				err = s3.DeleteImage(scheduledPost.Files[0].String())
				if err != nil {
					return
				}
				err = repository.ArchivizePost(scheduledPost.ID)
				if err != nil {
					return
				}
				notificationService.SendNotification("Notification", "Scheduled message sent!", scheduledPost.DeviceToken)
			case "mediaGroup":
				var files []*io.Reader
				var filenames []string
				for _, fileID := range scheduledPost.Files {
					media, err := s3.GetMedia(fileID.String())
					if err != nil {
						return
					}
					filenames = append(filenames, fileID.String())
					files = append(files, &media)
				}
				_, err := telegram.SendMediaGroup(files, filenames, scheduledPost.Text, scheduledPost.ChannelName)
				if err != nil {
					return
				}
				for _, fileID := range scheduledPost.Files {
					_ = s3.DeleteMedia(fileID.String())
				}
				err = repository.ArchivizePost(scheduledPost.ID)
				if err != nil {
					return
				}
				notificationService.SendNotification("Notification", "Scheduled message sent!", scheduledPost.DeviceToken)
			case "video":
				file, err := s3.GetVideo(scheduledPost.Files[0].String())
				if err != nil {
					return
				}
				_, err = telegram.SendVideoBytes(file, scheduledPost.Files[0].String(), scheduledPost.Text, scheduledPost.ChannelName)
				if err != nil {
					return
				}
				err = repository.ArchivizePost(scheduledPost.ID)
				notificationService.SendNotification("Notification", "Scheduled message sent!", scheduledPost.DeviceToken)
			case "audio":
				file, err := s3.GetAudio(scheduledPost.Files[0].String())
				if err != nil {
					return
				}
				_, err = telegram.SendAudioBytes(file, scheduledPost.Text, scheduledPost.ChannelName, scheduledPost.Files[0].String())
				if err != nil {
					return
				}
				err = repository.ArchivizePost(scheduledPost.ID)
				notificationService.SendNotification("Notification", "Scheduled message sent!", scheduledPost.DeviceToken)
			case "voice":
				file, err := s3.GetAudio(scheduledPost.Files[0].String())
				if err != nil {
					return
				}
				_, err = telegram.SendVoiceBytes(file, scheduledPost.Text, scheduledPost.ChannelName, scheduledPost.Files[0].String())
				if err != nil {
					return
				}
				//err = repository.DeleteScheduledPostById(scheduledPost.ID)
				err = repository.ArchivizePost(scheduledPost.ID)
				notificationService.SendNotification("Notification", "Scheduled message sent!", scheduledPost.DeviceToken)

			}
		} else if !scheduledPost.Scheduled.After(currentTime) {
			return
		}
	}

}
