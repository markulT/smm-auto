package scrapper

import (
	"fmt"
	"time"
)

type Channel struct {
	Name        string
	ChannelHash string
}

type ScrapperTask struct{}

func RunAnalyticsTask(chRepo ChannelAuthorizerRepository, s *Scrapper) {
	time.Sleep(10 * time.Second)
	RunAnalytics(chRepo, *s)
	ticker := time.NewTicker(1 * time.Hour)
	for {
		select {
		case <-ticker.C:
			RunAnalytics(chRepo, *s)
		}
	}
}

func RunAnalytics(chRepo ChannelAuthorizerRepository, s Scrapper) {
	var channelList []*Channel
	channelList = chRepo.GetAllChannels()
	initTime := time.Now()
	for _, channel := range channelList {
		if channel.ChannelHash != "" {
			fmt.Println("scanning by channel id")
			s.CollectAvgViewsByChannelID(channel.ChannelHash, channel.Name)
		} else {
			fmt.Println("scanning channel name")
			hash, err := s.CollectAvgViews(channel.Name)
			if err != nil {
				continue
			}
			chRepo.SaveChannelHash(channel.Name, hash)
		}
	}
	elapsed := time.Since(initTime)
	fmt.Printf("The whole counting process took %s to run.\n", elapsed)
}
