package scrapper

import "fmt"

// DefaultRepo This default implementation of ChannelAuthorizerRepository should only be used for testing purposes
type DefaultRepo struct{}

func (r *DefaultRepo) SaveCode(c Code) error {
	fmt.Println(c.SubmitCode)
	return nil
}
func (r *DefaultRepo) SaveChannelAvgViews(ch string, v float64) error {
	fmt.Println(ch, v)
	return nil
}

func (r *DefaultRepo) GetAllChannels() []*Channel {
	chArray := make([]*Channel, 0)
	chArray = append(chArray, &Channel{Name: "@privatnamemarnya", ChannelHash: ""})
	chArray = append(chArray, &Channel{Name: "@it_techgen", ChannelHash: "@it_techgen"})

	chArray = append(chArray, &Channel{Name: "@zhytel_ukraine", ChannelHash: ""})
	chArray = append(chArray, &Channel{Name: "@kharkovchanee", ChannelHash: ""})
	chArray = append(chArray, &Channel{Name: "@kpszsu", ChannelHash: ""})

	chArray = append(chArray, &Channel{Name: "@lvivyany_news", ChannelHash: ""})
	chArray = append(chArray, &Channel{Name: "@oleksiihoncharenko", ChannelHash: ""})
	chArray = append(chArray, &Channel{Name: "@jolybells", ChannelHash: ""})

	return chArray
}

//func (r *DefaultRepo) GetAllChannels() []*Channel {
//	chArray := make([]*Channel, 0)
//	chArray = append(chArray, &Channel{Name: "@privatnamemarnya", ChannelHash: ""})
//	chArray = append(chArray, &Channel{Name: "@it_techgen", ChannelHash: "@it_techgen"})
//
//	chArray = append(chArray, &Channel{Name: "@it_techgen", ChannelHash: "@it_techgen"})
//	chArray = append(chArray, &Channel{Name: "@it_techgen", ChannelHash: "@it_techgen"})
//	chArray = append(chArray, &Channel{Name: "@it_techgen", ChannelHash: "@it_techgen"})
//
//	chArray = append(chArray, &Channel{Name: "@it_techgen", ChannelHash: "@it_techgen"})
//	chArray = append(chArray, &Channel{Name: "@it_techgen", ChannelHash: "@it_techgen"})
//	chArray = append(chArray, &Channel{Name: "@it_techgen", ChannelHash: "@it_techgen"})
//
//	return chArray
//}

func (r *DefaultRepo) SaveChannelHash(ch string, hash string) error {
	fmt.Println(ch + " " + hash)
	return nil
}
