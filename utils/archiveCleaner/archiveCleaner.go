package archiveCleaner

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"golearn/models"
	"golearn/repository"
	"golearn/utils"
	"golearn/utils/s3"
	"os"
	"strconv"
	"sync"
	"time"
)



func RunArchiveCleaner() {
	ticker := time.NewTicker(1*time.Minute)
	defer ticker.Stop()
	for {
		select {
		case<-ticker.C:
			cleanArchive()
		}
	}
}

func cleanArchive()  {
	nElements := int(getTableLength("posts"))
	batchSize, _ :=  strconv.Atoi(os.Getenv("batchSize"))
	var wg sync.WaitGroup
	for start:=0;start < nElements; start += batchSize {
		end := start + batchSize -1
		if end >= nElements {
			end = nElements - 1
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
	fileRepo := repository.NewFileRepo()
	fileService := s3.NewFileService(fileRepo)

	defer wg.Done()

	var archivedPosts []models.Post

	archivedPosts = *repository.GetScheduledPostRelations(context.Background(),start, end-start+1, true)
	for _, archivedPost := range archivedPosts {
		ago := time.Now().AddDate(0,0,-14)
		if archivedPost.Scheduled.Before(ago) {
			//repository.ArchivizePost(archivedPost.ID)
			repository.DeleteScheduledPostById(context.Background(),archivedPost.ID)
			fileService.DeleteManyByID(context.Background(), archivedPost.Files)
		}
	}

}