package analytics

import "golearn/repository"

type AnalyticsService interface {}

type DefaultAnalyticsService struct {
	AnalyticsRepo repository.AnalyticsRepo
}

func NewAnalyticsService(ar repository.AnalyticsRepo) AnalyticsService {
	return DefaultAnalyticsService{AnalyticsRepo: ar}
}
