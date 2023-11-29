package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
	"golearn/models"
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

func SendMessage(botToken string,text string, chatId string) error {
	url := "https://api.telegram.org/bot" + botToken + "/sendMessage"
	sendMessageRequest := SendMessageRequest{
		Text:                  "" + text + "",
		DisableWebPagePreview: false,
		DisableNotification:   false,
		ReplyToMessage:        "",
		ChatId:                chatId,
	}
	//payload := strings.NewReader("{\"text\":\"Хочу присоромити одну дамочку\",\"parse_mode\":\"Optional\",\"disable_web_page_preview\":false,\"disable_notification\":false,\"reply_to_message_id\":null,\"chat_id\":\"@smm_auto_test\"}")
	jsonData, err := json.Marshal(sendMessageRequest)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("User-Agent", "Telegram Bot SDK - (https://github.com/irazasyed/telegram-bot-sdk)")
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))
	return nil
}

func SendDice(botToken string) error {
	url := "https://api.telegram.org/bot" + botToken + "/sendDice"
	var reqBody struct {
		ChatId string `json:"chat_id"`
	}
	reqBody.ChatId = "@smm_auto_test"
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("User-Agent", "Telegram Bot SDK - (https://github.com/irazasyed/telegram-bot-sdk)")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	_, err = io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return nil
}

func SendMediaGroupLinks(botToken string,links []string, caption string) (message string, err error) {
	url := "https://api.telegram.org/bot" + botToken + "/sendMediaGroup"

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

func SendMediaGroup(botToken string,files []*io.Reader,filenames[]string,fileModels []models.File, caption string, channelName string) (string, error) {


	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	fmt.Println(bot)
	//defer bot.Close()
	if err != nil {
		return "", err
	}

	var mediaItems []telego.InputMedia

	for index,file := range files {
		fileModel := fileModels[index]
		switch fileModel.Type {
		case "photo":
			fmt.Println("processing this file photo")
			fmt.Println(file)
			customReader := &CustomReader{reader:*file,filename: filenames[index]}
			fmt.Println(customReader.filename)
			media := telegoutil.MediaPhoto(telego.InputFile{
				File: customReader,
			})
			fmt.Println(media)

			if index == 0 {
				media = media.WithCaption(caption)
			}
			mediaItems = append(mediaItems, media)
		case "video":
			fmt.Println("processing this file video")
			fmt.Println(file)
			//customReader := &CustomReader{reader:*file,filename: filenames[index]}
			//media := telegoutil.MediaVideo(telego.InputFile{
			//	File: customReader,
			//})
			//if index == 0 {
			//	media = media.WithCaption(caption)
			//}
			//mediaItems = append(mediaItems, media)
		case "audio":
			fmt.Println("processing this file audio")
			fmt.Println(file)
			//customReader := &CustomReader{reader:*file,filename: filenames[index]}
			//media := telegoutil.MediaAudio(telego.InputFile{
			//	File: customReader,
			//})
			//if index == 0 {
			//	media = media.WithCaption(caption)
			//}
			//mediaItems = append(mediaItems, media)
		}
	}
	mdGroup := telegoutil.MediaGroup(telegoutil.Username(channelName), mediaItems...)
	fmt.Println(mdGroup)
	_, err = bot.SendMediaGroup(mdGroup)
	if err != nil {
		return "", err
	}

	return "", nil
}

func SendPhoto(botToken string,src io.Reader, caption string, fileName string, channelName string) (message string, err error) {
	url := "https://api.telegram.org/bot" + botToken + "/sendPhoto"
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
	_ = writer.WriteField("chat_id", channelName)
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


func SendAudio(botToken string,file *multipart.FileHeader, caption string, chatId string) (message string, err error) {
	url := "https://api.telegram.org/bot" + botToken + "/sendAudio"
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
		return "", err
	}
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	//body, _ := io.ReadAll(res.Body)
	return "success", nil
}

func SendAudioBytes(botToken string,reader io.Reader, caption string,chatId string, filename string) (string, error) {
	url := "https://api.telegram.org/bot" + botToken + "/sendAudio"
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	imageField, err := writer.CreateFormFile("audio", filename)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(imageField, reader)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	_ = writer.WriteField("chat_id", chatId)
	_ = writer.WriteField("caption", caption)
	err = writer.Close()
	if err != nil {
		fmt.Println(err)
		return "", err
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
	//body, _ := io.ReadAll(res.Body)
	return "success", nil
}

func SendVoiceBytes(botToken string, reader io.Reader, caption string,chatId string, filename string) (string, error) {
	url := "https://api.telegram.org/bot" + botToken + "/sendVoice"
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	imageField, err := writer.CreateFormFile("voice", filename)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(imageField, reader)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	_ = writer.WriteField("chat_id", chatId)
	_ = writer.WriteField("caption", caption)
	err = writer.Close()
	if err != nil {
		return "", err
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
	//body, _ := io.ReadAll(res.Body)
	return "success", nil
}

func SendVoice(botToken string,file *multipart.FileHeader, caption string, chatId string) (message string, err error) {
	url := "https://api.telegram.org/bot" + botToken + "/sendVoice"
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

func SendVideo(botToken string,file *os.File, caption string, chatId string, filename string) (string, error) {
	url := "https://api.telegram.org/bot" + botToken + "/sendVideo"


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

func SendVideoBytes(botToken string,file io.Reader, filename string, caption string, chatId string) ( string, error) {
	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
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



func SendVideoNote(botToken string,file *multipart.FileHeader, chatId string) (message string, err error) {
	url := "https://api.telegram.org/bot" + botToken + "/sendVideoNote"
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

func SendLocation(botToken string,latitude string, longitude string, chatId string) error {
	url := "https://api.telegram.org/bot" + botToken + "/sendLocation"
	sendMessageRequest := SendLocationRequest{
		Latitude:             latitude,
		Longitude:            longitude,
		ChatId:               chatId,
		ProximityAlertRadius: "90000",
	}
	//payload := strings.NewReader("{\"text\":\"Хочу присоромити одну дамочку\",\"parse_mode\":\"Optional\",\"disable_web_page_preview\":false,\"disable_notification\":false,\"reply_to_message_id\":null,\"chat_id\":\"@smm_auto_test\"}")
	jsonData, err := json.Marshal(sendMessageRequest)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("User-Agent", "Telegram Bot SDK - (https://github.com/irazasyed/telegram-bot-sdk)")
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	_, err = io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return nil
}

func SendVenue(botToken string, latitude string, longitude string, title string, address string, chatId string) error {
	url := "https://api.telegram.org/bot" + botToken + "/sendVenue"
	sendMessageRequest := SendVenueRequest{
		Latitude:     latitude,
		Longitude:    longitude,
		Title:        title,
		Address:      address,
		FoursquareId: "16015",
		ChatId:       chatId,
	}
	jsonData, err := json.Marshal(sendMessageRequest)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("User-Agent", "Telegram Bot SDK - (https://github.com/irazasyed/telegram-bot-sdk)")
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	_, err = io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return nil

}

