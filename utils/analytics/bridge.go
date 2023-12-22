package analytics

import (
	"context"
	"golearn/repository"
	"golearn/utils/scrapper"
)

type BridgeAnalyticsRepo struct {
	originalRepo repository.AnalyticsRepo
	channelRepo repository.ChannelRepository
}

func CreateBridgeRepo(repo repository.AnalyticsRepo) AnalyticsTaskRepo {
	return &BridgeAnalyticsRepo{originalRepo: repo}
}

func (br *BridgeAnalyticsRepo) SaveChannelAvgViews(ch string, v float64) error {
	return br.originalRepo.SaveChannelAvgViews(context.Background(),ch, v)
}

func (br *BridgeAnalyticsRepo) GetAllChannels() []*scrapper.Channel {
	var channels []*scrapper.Channel
	originalChannels, err := br.originalRepo.GetAllOriginalChannels(context.Background())
	if err != nil {
		return nil
	}

	for _,channel := range originalChannels {
		ch := &scrapper.Channel{
			Name:        channel.Name,
			ChannelHash: channel.ChannelHash,
		}
		channels = append(channels, ch)
	}

	return channels
}

func (br *BridgeAnalyticsRepo) SaveChannelHash(ch string, hash string) error {
	return br.channelRepo.SaveChannelHash(context.Background(), ch, hash)
}

