package analytics

import (
	"context"
	"golearn/utils/scrapper"
)

type AnalyticsTaskRepo interface {
	SaveChannelAvgViews(ch string, v float64) error
	GetAllChannels() []*scrapper.Channel
	SaveChannelHash(ch string, hash string) error
}

type AnalyticsLogService interface {
	LogError(s scrapper.ScrapperError)
}

type AnalyticsTask struct {
	Repo AnalyticsTaskRepo
	LogService AnalyticsLogService
}

func (at *AnalyticsTask) SetRepo(repo AnalyticsTaskRepo) {
	at.Repo = repo
}

func (at *AnalyticsTask) RunAnalyticsModule() {
	telegramScrapper := scrapper.NewDefaultScrapper(at.Repo, at.LogService)
	respch := make(chan struct{}, 1)
	ctx := context.Background()
	telegramScrapper.AuthorizeBot(ctx, respch)
	for {
		select {
		case <-respch:
			go scrapper.RunAnalyticsTask(at.Repo, &telegramScrapper)
		case <-ctx.Done():
		}
	}
}
