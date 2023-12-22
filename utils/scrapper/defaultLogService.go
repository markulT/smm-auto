package scrapper

import "fmt"

type ScrapperLogService struct{}

func (ls *ScrapperLogService) LogError(s ScrapperError) {
	fmt.Println(s.Data)
}
