package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func SendVideo(file *multipart.FileHeader, caption string, chatId string) (message string, err error) {
	url := "https://api.telegram.org/bot" + os.Getenv("botToken") + "/sendVideo"
	fmt.Println(url)
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	imageField, err := writer.CreateFormFile("video", file.Filename)
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

//func SendVideoNote(videoData []byte, chatId string) (message string, err error) {
//	url := "https://api.telegram.org/bot" + os.Getenv("botToken") + "/sendVideoNote"
//	fmt.Println(url)
//
//	var requestBody bytes.Buffer
//	writer := multipart.NewWriter(&requestBody)
//	videoField, err := writer.CreateFormFile("video_note", "video_note.mp4") // Provide a filename
//	if err != nil {
//		fmt.Println(err)
//		return "", err
//	}
//	_, err = videoField.Write(videoData)
//	if err != nil {
//		fmt.Println(err)
//		return "", err
//	}
//	_ = writer.WriteField("chat_id", chatId)
//
//	err = writer.Close()
//	if err != nil {
//		fmt.Println(err)
//		return "", err
//	}
//
//	req, err := http.NewRequest("POST", url, &requestBody)
//	if err != nil {
//		fmt.Println(err)
//		return "", err
//	}
//	req.Header.Set("Content-Type", writer.FormDataContentType())
//	res, err := http.DefaultClient.Do(req)
//	if err != nil {
//		fmt.Println(err)
//		return "", err
//	}
//	defer res.Body.Close()
//
//	body, _ := io.ReadAll(res.Body)
//	fmt.Println("response body : ")
//	fmt.Println(string(body))
//	return "success", nil
//}

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
