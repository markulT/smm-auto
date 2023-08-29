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
	Text string `json:"text"`
	ParseMode string `json:"parse_mode"`
	DisableWebPagePreview bool `json:"disable_web_page_preview"`
	DisableNotification bool `json:"disable_notification"`
	ReplyToMessage string `json:"reply_to_message"`
	ChatId string `json:"chat_id"`
}

func SendMessage(text string) {
	url := "https://api.telegram.org/bot" + os.Getenv("botToken") + "/sendMessage"
	sendMessageRequest := SendMessageRequest{
		Text:                  "" + text + "",
		DisableWebPagePreview: false,
		DisableNotification:   false,
		ReplyToMessage:        "",
		ChatId:                "@smm_auto_test",
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

func SendMediaGroupLinks(links []string,caption string) (message string, err error) {
	fmt.Println("made it here")
	url := "https://api.telegram.org/bot" + os.Getenv("botToken") + "/sendMediaGroup"

	var media []map[string]interface{}

	for _, link := range links {
		// Create a media object for the image
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

func SendMediaGroup(files []*multipart.FileHeader, caption string) (message string, err error) {
	url := "https://api.telegram.org/bot" + os.Getenv("botToken") + "/sendMediaGroup"
	fmt.Println(url)
	var media []map[string]interface{}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close()


	for _, file := range files {
		imageField, err := writer.CreateFormFile(file.Filename, file.Filename)
		if err != nil {
			return "", err
		}
		src, err := file.Open()
		if err != nil {
			return "", err
		}
		defer src.Close()
		_, err = io.Copy(imageField, src)
		if err != nil {
			return "", err
		}
	}

	for _, file := range files {
		// Create a media object for the image
		fmt.Println(file.Filename)
		media = append(media, map[string]interface{}{
			"type":  "photo",
			"media": "attach://" + file.Filename,
		})
	}
	//requestBody := map[string]interface{}{
	//	"chat_id": "@smm_auto_test",
	//	"media":   media,
	//}

	jsonBody, err := json.Marshal(media)
	fmt.Println(string(jsonBody))
	writer.WriteField("media", string(jsonBody))
	writer.WriteField("chat_id", "@smm_auto_test")
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	fmt.Println("res will be here")
	res, err := http.DefaultClient.Do(req)
	fmt.Println("response is here")
	if err!=nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	resBody, _ := io.ReadAll(res.Body)
	fmt.Println("response body : ")
	fmt.Println(string(resBody))
	return "success", nil
}

func SendPhoto(file *multipart.FileHeader, caption string) (message string, err error) {
	url := "https://api.telegram.org/bot" + os.Getenv("botToken") + "/sendPhoto"
	fmt.Println(url)
	src, err := file.Open()
	if err!=nil {
		return "", err
	}
	defer src.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	imageField, err := writer.CreateFormFile("photo",file.Filename)
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
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	fmt.Println("response body : ")
	fmt.Println(string(body))
	return "success", nil
}

func SendPhotoTest()  {
	url := "https://api.telegram.org/bot" + os.Getenv("botToken") + "/sendPhoto"
	imageFile, err := os.Open("public/ryan.jpg")
	if err!=nil {
		fmt.Println("Error opening image : ", err)
		return
	}
	defer imageFile.Close()
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	imageField, err := writer.CreateFormFile("photo", "ryan.jpg")
	_, err = io.Copy(imageField, imageFile)
	if err != nil {
		fmt.Println("Error copying image data:", err)
		return
	}
	_ = writer.WriteField("chat_id", "@smm_auto_test")
	err = writer.Close()
	if err !=nil {
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
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	fmt.Println("response body : ")
	fmt.Println(string(body))
}