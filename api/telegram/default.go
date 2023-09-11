package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)


func GetMe() {
	url := "https://api.telegram.org/bot" + os.Getenv("botToken") + "/getMe"
	req, _ := http.NewRequest("POST", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("User-Agent", "Telegram Bot SDK - (https://github.com/irazasyed/telegram-bot-sdk)")

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(string(body))
}

type SendMessageRequest struct {
	Text                  string `json:"text"`
	ParseMode             string `json:"parse_mode"`
	DisableWebPagePreview bool   `json:"disable_web_page_preview"`
	DisableNotification   bool   `json:"disable_notification"`
	ReplyToMessage        string `json:"reply_to_message"`
	ChatId                string `json:"chat_id"`
}

type SendLocationRequest struct {
	Latitude             string `json:"latitude"`
	Longitude            string `json:"longitude"`
	ChatId               string `json:"chat_id"`
	ProximityAlertRadius string `json:"proximity_alert_radius"`
}

type SendVenueRequest struct {
	Latitude     string `json:"latitude"`
	Longitude    string `json:"longitude"`
	Title        string `json:"title"`
	Address      string `json:"address"`
	FoursquareId string `json:"foursquare_id"`
	ChatId       string `json:"chat_id"`
}

type SendAudioMessageRequest struct {
	Caption string                `json:"caption"`
	Audio   *multipart.FileHeader `json:"audio"`
	ChatId  string                `json:"chat_id"`
}

func SendMessage(text string, chatId string) {
	url := "https://api.telegram.org/bot" + os.Getenv("botToken") + "/sendMessage"
	sendMessageRequest := SendMessageRequest{
		Text:                  "" + text + "",
		DisableWebPagePreview: false,
		DisableNotification:   false,
		ReplyToMessage:        "",
		ChatId:                chatId,
	}
	//payload := strings.NewReader("{\"text\":\"Хочу присоромити одну дамочку\",\"parse_mode\":\"Optional\",\"disable_web_page_preview\":false,\"disable_notification\":false,\"reply_to_message_id\":null,\"chat_id\":\"@smm_auto_test\"}")
	jsonData, _ := json.Marshal(sendMessageRequest)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))

	req.Header.Add("accept", "application/json")
	req.Header.Add("User-Agent", "Telegram Bot SDK - (https://github.com/irazasyed/telegram-bot-sdk)")
	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(string(body))
}

func SendDice()  {
	url := "https://api.telegram.org/bot" + os.Getenv("botToken") + "/sendDice"
	var reqBody struct {
		ChatId string `json:"chat_id"`
	}
	reqBody.ChatId = "@smm_auto_test"
	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Add("accept", "application/json")
	req.Header.Add("User-Agent", "Telegram Bot SDK - (https://github.com/irazasyed/telegram-bot-sdk)")
	req.Header.Add("content-type", "application/json")
	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(string(body))
}

func SendMediaGroupLinks(links []string, caption string) (message string, err error) {
	url := "https://api.telegram.org/bot" + os.Getenv("botToken") + "/sendMediaGroup"

	var media []map[string]interface{}

	for _, link := range links {
		fmt.Println(link)
		media = append(media, map[string]interface{}{
			"type":  "photo",
			"media": link,
		})
	}

	requestBody := map[string]interface{}{
		"chat_id": "@smm_auto_test",
		"media":   media,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))

	req.Header.Add("accept", "application/json")
	req.Header.Add("User-Agent", "Telegram Bot SDK - (https://github.com/irazasyed/telegram-bot-sdk)")
	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	resBody, _ := io.ReadAll(res.Body)

	fmt.Println(string(resBody))


	return "success", nil
}


type CustomReader struct {
	reader io.Reader
	filename string
}

func (cr *CustomReader) Read(p []byte) (n int, err error) {
	return cr.reader.Read(p)
}


func (cnr *CustomReader) Name() string {
	return cnr.filename
}

func SendMediaGroup(files []*io.Reader,filenames[]string, caption string) (string, error) {
	bot, err := telego.NewBot(os.Getenv("botToken"), telego.WithDefaultDebugLogger())
	if err != nil {
		return "", err
	}
	var mediaItems []telego.InputMedia
	for index,file := range files {
		customReader := &CustomReader{reader:*file,filename: filenames[index]}
		media := telegoutil.MediaPhoto(telego.InputFile{
			File: customReader,
		})
		if index == 0 {
			media = media.WithCaption(caption)
		}
		mediaItems = append(mediaItems, media)
	}
	mdGroup := telegoutil.MediaGroup(telegoutil.Username("@smm_auto_test"), mediaItems...)
	_, _ = bot.SendMediaGroup(mdGroup)
	return "", nil
}

func SendPhoto(src io.Reader, caption string, fileName string) (message string, err error) {
	url := "https://api.telegram.org/bot" + os.Getenv("botToken") + "/sendPhoto"
	fmt.Println(url)

	if err!=nil {
		return "", err
	}


	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	imageField, err := writer.CreateFormFile("photo",fileName)
	if err!=nil {
		fmt.Println(err)
		return "", err
	}
	_, err = io.Copy(imageField, src)
	if err!=nil {
		fmt.Println(err)
		return "", nil
	}
	_ = writer.WriteField("chat_id", "@smm_auto_test")
	_ = writer.WriteField("caption", caption)
	err = writer.Close()
	if err !=nil {
		fmt.Println(err)
		return
	}
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := http.DefaultClient.Do(req)
	if err!=nil {
		fmt.Println(err)
		return
	}
	fmt.Println("aboba")
	fmt.Println(res)
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	fmt.Println("response body : ")
	fmt.Println(string(body))
	return "success", nil
}


func SendAudio(file *multipart.FileHeader, caption string, chatId string) (message string, err error) {
	url := "https://api.telegram.org/bot" + os.Getenv("botToken") + "/sendAudio"
	fmt.Println(url)
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	imageField, err := writer.CreateFormFile("audio", file.Filename)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	_, err = io.Copy(imageField, src)
	if err != nil {
		fmt.Println(err)
		return "", nil
	}
	_ = writer.WriteField("chat_id", chatId)
	_ = writer.WriteField("caption", caption)
	err = writer.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	fmt.Println("response body : ")
	fmt.Println(string(body))
	return "success", nil
}



func SendVoice(file *multipart.FileHeader, caption string, chatId string) (message string, err error) {
	url := "https://api.telegram.org/bot" + os.Getenv("botToken") + "/sendVoice"
	fmt.Println(url)
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	imageField, err := writer.CreateFormFile("voice", file.Filename)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	_, err = io.Copy(imageField, src)
	if err != nil {
		fmt.Println(err)
		return "", nil
	}
	_ = writer.WriteField("chat_id", chatId)
	_ = writer.WriteField("caption", caption)
	err = writer.Close()
	if err != nil {

		fmt.Println(err)
		return
	}
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {

		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := http.DefaultClient.Do(req)
	if err!=nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	fmt.Println("response body : ")
	fmt.Println(string(body))

	return "success", nil
}

func SendVideo(file *os.File, caption string, chatId string, filename string) (string, error) {
	fmt.Println("sending video")
	url := "https://api.telegram.org/bot" + os.Getenv("botToken") + "/sendVideo"


	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	imageField, err := writer.CreateFormFile("video", filename)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(imageField, file)
	if err != nil {
		fmt.Println(
			"error copying huynia")
		fmt.Println(err)
		return "", err
	}
	_ = writer.WriteField("chat_id", chatId)
	_ = writer.WriteField("caption", caption)
	err = writer.Close()
	if err != nil {
		return  "",err
	}
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	fmt.Println("response body : ")
	fmt.Println(string(body))
	return "success", nil
}
//func SendVideoBytes(file io.Reader, filename string, caption string, chatId string) (message string, err error) {
//
//	url := "https://api.telegram.org/bot" + os.Getenv("botToken") + "/sendVideo"
//	var requestBody bytes.Buffer
//	writer := multipart.NewWriter(&requestBody)
//
//	imageField, err := writer.CreateFormFile("video", filename)
//	if err != nil {
//		return "", err
//	}
//	fmt.Println(file == nil)
//	_, err = io.Copy(imageField, file)
//	if err != nil {
//		fmt.Println("error sending 1")
//		fmt.Println(err.Error())
//		return "", err
//	}
//
//	err = writer.WriteField("chat_id", chatId)
//	if err != nil {
//		return "", err
//	}
//	err = writer.WriteField("caption", caption)
//	if err != nil {
//		return "", err
//	}
//	if err != nil {
//		fmt.Println("error sending 0")
//		return "", err
//	}
//
//
//	defer writer.Close()
//	req, err := http.NewRequest("POST", url, &requestBody)
//	if err != nil {
//		fmt.Println("error sending 2")
//		return "", err
//	}
//	req.Header.Set("Content-Type", writer.FormDataContentType())
//	res, err := http.DefaultClient.Do(req)
//	if err != nil {
//		fmt.Println("error sending 3")
//		return
//	}
//	defer res.Body.Close()
//	body, _ := io.ReadAll(res.Body)
//	fmt.Println("response body : ")
//	fmt.Println(string(body))
//	return "success", nil
//
//}

func SendVideoBytes(file io.Reader, filename string, caption string, chatId string) ( string, error) {
	bot, err := telego.NewBot(os.Getenv("botToken"), telego.WithDefaultDebugLogger())
	if err != nil {
		return "",err
	}
	customReader := &CustomReader{
		reader:   file,
		filename: filename,
	}
	video := &telego.InputFile{
		File: customReader,
	}
	_, err = bot.SendVideo(&telego.SendVideoParams{
		ChatID:                   telego.ChatID{
			Username: chatId,
		},
		Video:                    *video,
		Caption:                  caption,
	})

	if err != nil {
		return "",err
	}
	return "", nil
}



func SendVideoNote(file *multipart.FileHeader, chatId string) (message string, err error) {
	url := "https://api.telegram.org/bot" + os.Getenv("botToken") + "/sendVideoNote"
	fmt.Println(url)
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	imageField, err := writer.CreateFormFile("video_note", file.Filename)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	_, err = io.Copy(imageField, src)
	if err != nil {
		fmt.Println(err)
		return "", nil
	}
	_ = writer.WriteField("chat_id", chatId)

	err = writer.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	fmt.Println("response body : ")
	fmt.Println(string(body))
	return "success", nil
}

func SendLocation(latitude string, longitude string, chatId string) {
	url := "https://api.telegram.org/bot" + os.Getenv("botToken") + "/sendLocation"
	sendMessageRequest := SendLocationRequest{
		Latitude:             latitude,
		Longitude:            longitude,
		ChatId:               chatId,
		ProximityAlertRadius: "90000",
	}
	//payload := strings.NewReader("{\"text\":\"Хочу присоромити одну дамочку\",\"parse_mode\":\"Optional\",\"disable_web_page_preview\":false,\"disable_notification\":false,\"reply_to_message_id\":null,\"chat_id\":\"@smm_auto_test\"}")
	jsonData, _ := json.Marshal(sendMessageRequest)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))

	req.Header.Add("accept", "application/json")
	req.Header.Add("User-Agent", "Telegram Bot SDK - (https://github.com/irazasyed/telegram-bot-sdk)")
	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(string(body))
}

func SendVenue(latitude string, longitude string, title string, address string, chatId string) {
	url := "https://api.telegram.org/bot" + os.Getenv("botToken") + "/sendVenue"
	sendMessageRequest := SendVenueRequest{
		Latitude:     latitude,
		Longitude:    longitude,
		Title:        title,
		Address:      address,
		FoursquareId: "16015",
		ChatId:       chatId,
	}
	jsonData, _ := json.Marshal(sendMessageRequest)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))

	req.Header.Add("accept", "application/json")
	req.Header.Add("User-Agent", "Telegram Bot SDK - (https://github.com/irazasyed/telegram-bot-sdk)")
	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(string(body))
}

