package scrapper

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

func RunSubmittionCodesCleaner() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:

		}
	}
}
func RunCleaning() {
	numElements := int(getTableLength("posts"))
	batchSize, _ := strconv.Atoi(os.Getenv("batchSize"))
	var wg *sync.WaitGroup
	for start := 0; start < numElements; start += batchSize {
		end := start + batchSize - 1
		if end >= numElements {
			end = numElements - 1
		}
		wg.Add(1)
		go processCodesBatch(start, end, wg)
	}
	wg.Wait()
}

func getTableLength(tableName string) int64 {
	return int64(rand.Int())
}

func processCodesBatch(start, end int, wg *sync.WaitGroup) {
	defer wg.Done()

	var codesBatch []Code
	// TODO: Add codes batches receiving function

	//
	for _, code := range codesBatch {
		fmt.Println(code)
	}
}
