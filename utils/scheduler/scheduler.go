package scheduler

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"golearn/api/telegram"
	"golearn/models"
	"golearn/repository"
	"golearn/utils"
	"os"
	"strconv"
	"sync"
	"time"
)

func FetchAndProcessPosts()  {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
			case <- ticker.C:
				fmt.Println("tick")
				ProcessPosts()
		}
	}
}

func ProcessPosts()  {
	numElements := int(getTableLength("posts"))
	batchSize, _ :=  strconv.Atoi(os.Getenv("batchSize"))
	fmt.Println(numElements)
	fmt.Println(batchSize)
	var wg sync.WaitGroup
	for start:=0;start < numElements; start += batchSize {
		end := start + batchSize -1
		if end >= numElements {
			end = numElements - 1
		}
		wg.Add(1)
		go processBatch(start, end, &wg)
	}
	wg.Wait()
}

func getTableLength(tableName string) int64 {
	collection:=utils.DB.Collection(tableName)

	count, err := collection.CountDocuments(context.Background(), bson.D{})
	if err != nil {
		return 0
	}
	return count
}

func processBatch(start, end int, wg *sync.WaitGroup)  {

	fmt.Println(start, end)

	defer wg.Done()

	var scheduledPosts []models.Post

	scheduledPosts = *repository.GetScheduledPostRelations(start, end-start+1)

	for _, scheduledPost := range scheduledPosts {
		originalTimezone, offset := scheduledPost.Scheduled.Zone()
		currentTime := time.Now().In(time.FixedZone(originalTimezone, offset))
		fmt.Println(scheduledPost.Scheduled)
		fmt.Println(currentTime)
		fmt.Println(scheduledPost.Scheduled.Before(currentTime))
		if scheduledPost.Scheduled.Before(currentTime) {
			switch scheduledPost.Type {
				case "message":
					telegram.SendMessage(scheduledPost.Text, scheduledPost.ChannelName)
					fmt.Println(scheduledPost.ID)
					err := repository.DeleteScheduledPostById(scheduledPost.ID)
					if err != nil {
						return
					}
				case "photo":

					telegram.SendPhoto()

			}
		} else if !scheduledPost.Scheduled.After(currentTime) {

		}
	}

}