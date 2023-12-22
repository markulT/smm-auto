package scrapper

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Scrapper interface {
	AuthorizeBot(c context.Context, respch chan struct{})
	CollectAvgViews(channelName string) (string, error)
	CollectAvgViewsByChannelID(channelID string, chn string) error
}

type ChannelAuthorizerRepository interface {
	SaveChannelAvgViews(ch string, v float64) error
	GetAllChannels() []*Channel
	SaveChannelHash(ch string, hash string) error
}

type ScrapperError struct {
	Data string `bson:"data" json:"data"`
}

type LogService interface {
	LogError(s ScrapperError)
}

type defaultScrapperImpl struct {
	ScrapperContext context.Context
	ChRepo          ChannelAuthorizerRepository
	LogService      LogService
}

func NewDefaultScrapper(chRepo ChannelAuthorizerRepository, ls LogService) Scrapper {
	scrapperContext, _ := chromedp.NewContext(context.Background())
	return &defaultScrapperImpl{ScrapperContext: scrapperContext, ChRepo: chRepo, LogService: ls}
}

func getChannelIDFromURL(url string) (string, error) {

	var numberValue string

	re := regexp.MustCompile(`#(-?\d+)`)
	match := re.FindStringSubmatch(url)

	if len(match) > 1 {
		numberValue = match[1]
		if numberValue[0] == '-' {
			numberValue = numberValue[1:]
		}
	} else {
	}
	return numberValue, nil
}

func (s *defaultScrapperImpl) CollectAvgViewsByChannelID(channelID string, channelName string) error {
	startTime := time.Now()

	var err error
	fmt.Println("https://web.telegram.org/k/" + "#" + channelID)
	// Go to channel using his telegram ID AKA hash
	err = chromedp.Run(s.ScrapperContext,
		chromedp.Navigate("https://web.telegram.org/k/"+"#"+channelID),
		chromedp.WaitReady("body"),
		//chromedp.WaitReady(".ChatInfo>.info>.title>.fullName", chromedp.ByQuery),
	)
	if err != nil {
		s.LogService.LogError(ScrapperError{Data: err.Error()})
		return err
	}
	fmt.Println("A1")
	var originalText string
	if originalText != "" {
		originalText = ""
	}
	//
	contextTest, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	pollRespCh := make(chan struct{}, 1)
	go func(ch chan struct{}) {
		fmt.Println("pojebano")
		err = chromedp.Run(s.ScrapperContext, chromedp.PollFunction(`() => {
	       	var header = document.querySelectorAll(".chat-info>.person>.content>.top>.user-title>.peer-title")[0];
			if (header != undefined) {	
	       		return header?.innerText != '`+originalText+`';
			} else {
				return false;
			}
	   }`, nil))
		if err == nil {
			fmt.Println("piszemy w porozny struct")
			pollRespCh <- struct{}{}
		} else {
			fmt.Println(err.Error())
			fmt.Println("pojebano 2")
		}
	}(pollRespCh)
	fmt.Println("przed cholera")
	select {
	case <-pollRespCh:
		fmt.Println("nic nie robie od teraz")
	case <-contextTest.Done():
		fmt.Println("done")
		var htmlContent string
		err = chromedp.Run(s.ScrapperContext,
			chromedp.OuterHTML("html", &htmlContent, chromedp.ByQuery),
		)
		//fmt.Println(htmlContent)
		var currentURL string
		err = chromedp.Run(s.ScrapperContext, chromedp.Location(&currentURL))
		//fmt.Println(currentURL)
	}
	fmt.Println("cholera zaczyna sie")
	cancel()
	fmt.Println("canceled")
	//var fullName string
	//err = chromedp.Run(s.ScrapperContext, chromedp.Text(`.ChatInfo>.info>.title>.fullName`, &fullName, chromedp.ByQuery))
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return err
	//}
	//fmt.Println(fullName)

	fmt.Println("A2")

	err = chromedp.Run(s.ScrapperContext, chromedp.WaitVisible("div.message", chromedp.ByQuery))
	fmt.Println("A3")

	fmt.Println("A4")

	//err = chromedp.Run(s.ScrapperContext, chromedp.Text(`.ChatInfo>.info>.title>.fullName`, &s.PrevChannel))
	fmt.Println("A5")
	var postsCount int
	err = chromedp.Run(s.ScrapperContext, chromedp.Evaluate(`
		(function(){
			let arr = Array.from(document.querySelectorAll(".post-views"))
			return arr.length
		})()
	`, &postsCount))
	fmt.Println("post count on page: ", postsCount)

	var avgViews int
	err = chromedp.Run(s.ScrapperContext, chromedp.Evaluate(`
           (function() {
   			let arr = Array.from(document.querySelectorAll(".post-views")).map(el => el.innerText).filter(str=>str!=='');
   			let numbers = arr.map(str => {
       			if (str.endsWith('K')) {
           			return parseFloat(str.slice(0, -1)) * 1000;
       			} else {
           			return parseInt(str);
       			}
   			});
   			return Math.floor(numbers.reduce((a, b) => a + b, 0) / numbers.length);
})()
	`, &avgViews))
	fmt.Println(avgViews)

	if err != nil {
		s.LogService.LogError(ScrapperError{Data: err.Error()})
		return err
	}
	fmt.Println("A6")

	s.ChRepo.SaveChannelAvgViews(channelName, float64(avgViews))

	if err != nil {
		return err
	}
	fmt.Printf("Time spent to count single channel: %v \n", time.Since(startTime).Seconds())
	return nil
}

func (s *defaultScrapperImpl) AuthorizeBot(c context.Context, respch chan struct{}) {

	var htmlContent string

	// 1 - open browser
	err := chromedp.Run(s.ScrapperContext,
		chromedp.Navigate("https://web.telegram.org/a"),
	)
	if err != nil {
		fmt.Println(err)
		c.Done()
	}

	err = chromedp.Run(s.ScrapperContext,
		chromedp.WaitVisible(`//button[normalize-space(text())="Log in by phone Number"]`),
	)

	if err != nil {
		log.Fatal(err)
	}

	err = chromedp.Run(s.ScrapperContext,
		chromedp.Click(`//button[normalize-space(text())="Log in by phone Number"]`),
	)

	if err != nil {
		log.Fatal(err)
	}

	err = chromedp.Run(s.ScrapperContext,
		chromedp.WaitVisible(`input[id="sign-in-phone-number"]`),
	)

	if err != nil {
		log.Fatal(err)
	}

	err = chromedp.Run(s.ScrapperContext,
		chromedp.SendKeys(`input[id="sign-in-phone-number"]`, "+380 98 799 74 10"),
		//chromedp.SendKeys(`input[id="sign-in-phone-number"]`, "+380 67 766 55 69"),
		//chromedp.SendKeys(`#sign-in-phone-number`, "+380 98 799 74 10", chromedp.ByQuery),
	)

	if err != nil {
		fmt.Println(err)
		c.Done()
	}

	var inputTestValue string
	err = chromedp.Run(s.ScrapperContext, chromedp.EvaluateAsDevTools(`document.getElementById("sign-in-phone-number").value`, &inputTestValue))
	if err != nil {
		panic(err)
	}
	fmt.Println(inputTestValue)

	var value string
	err = chromedp.Run(s.ScrapperContext, chromedp.Text(`label[for="sign-in-phone-number"]`, &value, chromedp.ByQuery))
	if err != nil {
		panic(err)
	}
	time.Sleep(2 * time.Second)
	err = chromedp.Run(s.ScrapperContext,
		chromedp.EvaluateAsDevTools(`document.querySelectorAll(".Button")[0].click()`, nil),
	)

	if err != nil {
		fmt.Println(err)
		c.Done()
	}

	time.Sleep(5 * time.Second)

	var newValue string
	err = chromedp.Run(s.ScrapperContext, chromedp.Text(`label[for="sign-in-phone-number"]`, &newValue, chromedp.ByQuery))
	if err != nil {
		panic(err)
	}

	err = chromedp.Run(s.ScrapperContext,
		chromedp.OuterHTML("html", &htmlContent, chromedp.ByQuery),
	)
	if err != nil {
		log.Fatal(err)
	}

	// 3 - wait...

	err = chromedp.Run(s.ScrapperContext, chromedp.WaitVisible(`#sign-in-code`))
	if err != nil {
		c.Done()
	}

	var activationCode string
	fmt.Println("Input activation code : ")
	fmt.Scanln(&activationCode)

	// 4 - code inputted

	err = chromedp.Run(s.ScrapperContext,
		chromedp.WaitVisible(`input#sign-in-code`, chromedp.ByQuery),
	)
	if err != nil {
		c.Done()
	}

	err = chromedp.Run(s.ScrapperContext,
		chromedp.SendKeys(`input[id="sign-in-code"]`, activationCode, chromedp.NodeVisible),
	)
	if err != nil {
		c.Done()
	}

	respch <- struct{}{}

}

func calcAvgViews(n []string) float64 {
	total := 0.0
	for _, item := range n {
		if strings.Contains(item, "K") {
			number, _ := strconv.ParseFloat(strings.Replace(item, "K", "", -1), 64)
			total += number * 1000
		} else {
			number, _ := strconv.ParseFloat(item, 64)
			total += number
		}
	}
	return total / float64(len(n))
}

//func (s *defaultScrapperImpl) CollectAvgViews(channelName string) error {
//	var err error
//
//	return nil
//}

func (s *defaultScrapperImpl) CollectAvgViews(channelName string) (string, error) {
	startTime := time.Now()
	var err error
	err = chromedp.Run(s.ScrapperContext,
		chromedp.Navigate("https://web.telegram.org/a"),
		chromedp.WaitReady("body"),
	)

	fmt.Println("1")
	err = chromedp.Run(s.ScrapperContext, chromedp.SendKeys(`input#telegram-search-input`, channelName))
	if err != nil {
		fmt.Println("")
		s.LogService.LogError(ScrapperError{Data: err.Error()})
		return "", err
	}
	fmt.Println("2")
	err = chromedp.Run(s.ScrapperContext, chromedp.WaitVisible(`.ListItem.search-result`))
	if err != nil {
		s.LogService.LogError(ScrapperError{Data: err.Error()})
		return "", err
	}
	var originalText string
	if originalText != "" {
		originalText = " "
	}
	fmt.Println("3")
	err = chromedp.Run(s.ScrapperContext, chromedp.Click(`.ListItem.search-result`, chromedp.NodeVisible))
	if err != nil {
		s.LogService.LogError(ScrapperError{Data: err.Error()})
		return "", err
	}

	fmt.Println("4")
	err = chromedp.Run(s.ScrapperContext, chromedp.PollFunction(`() => {
	       let header = document.querySelectorAll(".ChatInfo>.info>.title>.fullName")[0];
	       return header.innerText != `+"`"+originalText+"`"+`;
	   }`, nil))
	if err != nil {
		s.LogService.LogError(ScrapperError{Data: err.Error()})
		return "", err
	}
	fmt.Println("5")
	err = chromedp.Run(s.ScrapperContext, chromedp.Text(`.ChatInfo>.info>.title>.fullName`, &originalText))
	if err != nil {
		s.LogService.LogError(ScrapperError{Data: err.Error()})
		return "", err
	}
	fmt.Println("6")
	var htmlContent string
	err = chromedp.Run(s.ScrapperContext,
		chromedp.OuterHTML("html", &htmlContent, chromedp.ByQuery),
	)

	//contextTest, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//select {
	//case <-contextTest.Done():
	//	fmt.Println(htmlContent)
	//}
	err = chromedp.Run(s.ScrapperContext, chromedp.WaitVisible("div.Message", chromedp.ByQuery))
	if err != nil {

		fmt.Println(err.Error())
		s.LogService.LogError(ScrapperError{Data: err.Error()})
		return "", err
	}

	fmt.Println("7")
	err = chromedp.Run(s.ScrapperContext, chromedp.PollFunction(`() => {
	       let header = document.querySelectorAll("div.Message")[document.querySelectorAll("div.Message").length-1];
	       return header.innerText != `+"`"+originalText+"`"+`;
	   }`, nil))
	if err != nil {
		s.LogService.LogError(ScrapperError{Data: err.Error()})
		return "", err
	}
	fmt.Println("8")
	//err = chromedp.Run(s.ScrapperContext, chromedp.Text(`.ChatInfo>.info>.title>.fullName`, &s.PrevChannel))

	var postsCount int
	err = chromedp.Run(s.ScrapperContext, chromedp.Evaluate(`
		(function(){
			let arr = Array.from(document.querySelectorAll(".MessageMeta>.message-views"))
			return arr.length
		})()
	`, &postsCount))
	fmt.Println("9")
	var avgViews int
	err = chromedp.Run(s.ScrapperContext, chromedp.Evaluate(`
           (function() {
   			let arr = Array.from(document.querySelectorAll(".MessageMeta>.message-views")).map(el => el.innerText);
   			let numbers = arr.map(str => {
       			if (str.endsWith('K')) {
           			return parseFloat(str.slice(0, -1)) * 1000;
       			} else {
           			return parseInt(str);
       			}
   			});
   			return Math.floor(numbers.reduce((a, b) => a + b, 0) / numbers.length);
})()
	`, &avgViews))

	if err != nil {
		s.LogService.LogError(ScrapperError{
			Data: err.Error(),
		})
		return "", err
	}
	fmt.Println("10")
	err = s.ChRepo.SaveChannelAvgViews(channelName, float64(avgViews))
	if err != nil {
		s.LogService.LogError(ScrapperError{Data: err.Error()})
		return "", err
	}

	err = chromedp.Run(s.ScrapperContext,
		chromedp.Navigate("https://web.telegram.org/a"),
		chromedp.WaitReady("body"),
	)
	fmt.Printf("Time spent to count single channel: %v \n", time.Since(startTime).Seconds())

	return "", nil
}

//func checkForValueChange(ctx context.Context, sel string, prevValue string) chromedp.Action {
//	return chromedp.ActionFunc(func(ctx context.Context, h cdp.Executor) error {
//		var newValue string
//		err := chromedp.EvaluateAsDevTools(sel, &newValue).Do(ctx, h)
//		if err != nil {
//			return err
//		}
//		if newValue != prevValue {
//			return chromedp.ErrContinue
//		}
//		return nil
//	})
//}
