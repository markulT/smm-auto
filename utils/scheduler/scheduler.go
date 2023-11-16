package scheduler

import (
	"context"
	"firebase.google.com/go/messaging"
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

type SchedulerTask struct {
	FcmClient *messaging.Client
}

func (s *SchedulerTask) FetchAndProcessPosts()  {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <- ticker.C:
			s.ProcessPosts()
		}
	}
}

func (s *SchedulerTask) ProcessPosts()  {
	numElements := int(s.getTableLength("posts"))
	batchSize, _ :=  strconv.Atoi(os.Getenv("batchSize"))
	var wg sync.WaitGroup
	for start:=0;start < numElements; start += batchSize {
		end := start + batchSize -1
		if end >= numElements {
			end = numElements - 1
		}
		wg.Add(1)
		go s.processBatch(start, end, &wg)
	}
	wg.Wait()
}

func (s *SchedulerTask) getTableLength(tableName string) int64 {
	collection:=utils.DB.Collection(tableName)

	count, err := collection.CountDocuments(context.Background(), bson.D{})
	if err != nil {
		return 0
	}
	return count
}

func (s *SchedulerTask) processBatch(start, end int, wg *sync.WaitGroup)  {
	defer wg.Done()

	var scheduledPosts []models.Post

	notificationService := notifications.DefaultNotificationService{}

	scheduledPosts = *repository.GetScheduledPostRelations(context.Background(),start, end-start+1, false)

	for _, scheduledPost := range scheduledPosts {
		originalTimezone, offset := scheduledPost.Scheduled.Zone()
		currentTime := time.Now().In(time.FixedZone(originalTimezone, offset))
		if scheduledPost.Scheduled.Before(currentTime) {

			channelRepo := repository.NewChannelRepo()
			channel, err := channelRepo.FindByID(context.Background(), scheduledPost.ChannelID)
			if err != nil {
				continue
			}

			switch scheduledPost.Type {
			case "message":
				telegram.SendMessage(channel.AssignedBotToken,scheduledPost.Text, channel.Name)
				err := repository.ArchivizePost(context.Background(),scheduledPost.ID)
				if err != nil {
					return
				}
				notificationService.SendNotification("Notification", "Scheduled message sent!", scheduledPost.DeviceToken)
			case "photo":
				image, err := s3.GetImage(scheduledPost.Files[0].String())
				if err != nil {
					return
				}
				_, err = telegram.SendPhoto(channel.AssignedBotToken,image, scheduledPost.Text, scheduledPost.Files[0].String(), channel.Name)
				if err != nil {
					return
				}
				err = s3.DeleteImage(scheduledPost.Files[0].String())
				if err != nil {
					return
				}
				err = repository.ArchivizePost(context.Background(),scheduledPost.ID)
				if err != nil {
					return
				}
				notificationService.SendNotification("Notification", "Scheduled message sent!", scheduledPost.DeviceToken)
			case "mediaGroup":
				fileRepo := repository.NewFileRepo()
				var files []*io.Reader
				var filenames []string
				var fileModels []models.File
				for _, fileID := range scheduledPost.Files {
					media, err := s3.GetMedia(fileID.String())
					if err != nil {
						return
					}
					fileModel, err := fileRepo.FindByID(context.Background(), fileID)
					filenames = append(filenames, fileID.String())
					files = append(files, &media)
					fileModels = append(fileModels, *fileModel)
				}
				_, err := telegram.SendMediaGroup(channel.AssignedBotToken,files, filenames,fileModels, scheduledPost.Text, channel.Name)
				if err != nil {
					return
				}
				for _, fileID := range scheduledPost.Files {
					_ = s3.DeleteMedia(fileID.String())
				}
				err = repository.ArchivizePost(context.Background(),scheduledPost.ID)
				if err != nil {
					return
				}
				notificationService.SendNotification("Notification", "Scheduled message sent!", scheduledPost.DeviceToken)
			case "video":
				file, err := s3.GetVideo(scheduledPost.Files[0].String())
				if err != nil {
					return
				}
				_, err = telegram.SendVideoBytes(channel.AssignedBotToken,file, scheduledPost.Files[0].String(), scheduledPost.Text, channel.Name)
				if err != nil {
					return
				}
				err = repository.ArchivizePost(context.Background(),scheduledPost.ID)
				notificationService.SendNotification("Notification", "Scheduled message sent!", scheduledPost.DeviceToken)
			case "audio":
				file, err := s3.GetAudio(scheduledPost.Files[0].String())
				if err != nil {
					return
				}
				_,err = telegram.SendAudioBytes(channel.AssignedBotToken,file, scheduledPost.Text, channel.Name, scheduledPost.Files[0].String())
				if err != nil {
					return
				}
				err = repository.ArchivizePost(context.Background(),scheduledPost.ID)
				notificationService.SendNotification("Notification", "Scheduled message sent!", scheduledPost.DeviceToken)
			case "voice":
				file, err := s3.GetAudio(scheduledPost.Files[0].String())
				if err != nil {
					return
				}
				_,err = telegram.SendVoiceBytes(channel.AssignedBotToken,file, scheduledPost.Text, channel.Name, scheduledPost.Files[0].String())
				if err != nil {
					return
				}
				//err = repository.DeleteScheduledPostById(scheduledPost.ID)
				err = repository.ArchivizePost(context.Background(),scheduledPost.ID)
				notificationService.SendNotification("Notification", "Scheduled message sent!", scheduledPost.DeviceToken)

			}
		} else if !scheduledPost.Scheduled.After(currentTime) {
			return
		}
	}

}