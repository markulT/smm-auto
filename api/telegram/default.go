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

//func SendMediaGroup(files []*multipart.FileHeader, caption string) (string, error) {
//	requestUrl := "https://api.telegram.org/bot" + os.Getenv("botToken") + "/getMe"
//	var requestBody bytes.Buffer
//	writer:=multipart.NewWriter(&requestBody)
//	for index, file := range files {
//		src, err := file.Open()
//		if err!=nil {
//			return "", err
//		}
//		defer src.Close()
//		fileField, err := writer.CreateFormFile(fmt.Sprintf("image%d", index), file.Filename)
//		_, err = io.Copy(fileField, src)
//		if err != nil {
//			return "", err
//		}
//	}
//	writer.WriteField("chat_id", "@smm_auto_test")
//	writer.WriteField("media", `[{"type":"photo", "media":"attach://image0"}, {"type":"photo", "media":"attach://image1"}]`)
//	defer writer.Close()
//	req, err := http.NewRequest("POST", requestUrl, &requestBody)
//	if err != nil {
//		fmt.Println("Error creating request:", err)
//		return "", err
//	}
//	req.Header.Set("Content-Type", "multipart/form-data")
//	client := &http.Client{}
//	res, err := client.Do(req)
//	fmt.Println(res)
//	fmt.Println(res.Status)
//	fmt.Println(res.StatusCode)
//
//	if err != nil {
//		fmt.Println("Error sending request:", err)
//		return "", err
//	}
//	defer res.Body.Close()
//	return "a", nil
//}

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

func SendMediaGroupLinks(links []string,caption string) (message string, err error) {
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

type NamedReaderCloser struct {
	io.Reader
	multipart.FileHeader
}

func (c *NamedReaderCloser) Name() string {
	return c.FileHeader.Filename
}

type CustomReader struct {
	file multipart.File
}

func (cr *CustomReader) Read(p []byte) (n int, err error) {
	return cr.file.Read(p)
}

func SendMediaGroupLazy(files []*multipart.FileHeader, caption string) (string, error) {
	bot, err := telego.NewBot(os.Getenv("botToken"), telego.WithDefaultDebugLogger())
	if err != nil {
		return "", err
	}
	var mediaItems []telego.InputMedia
	for index, file := range files {
		src, err := file.Open()
		defer src.Close()
		if err != nil {
			return "", err
		}
		customReader := &CustomReader{file: src}
		closer := NamedReaderCloser{
			Reader: customReader,
			FileHeader: *file,
		}
		media := telegoutil.MediaPhoto(telego.InputFile{
			File: &closer,
		})
		if index == 0 {
			media = media.WithCaption(caption)
		}
		mediaItems = append(mediaItems, media)
	}
	mdGroup := telegoutil.MediaGroup(telegoutil.Username("@smm_auto_test"), mediaItems...)
	_, _ = bot.SendMediaGroup(mdGroup)
	return "success", nil
}

//func SendMediaGroup(files []*multipart.FileHeader, caption string) (string, error) {
//	url := "https://api.telegram.org/bot" + os.Getenv("botToken") + "/sendMediaGroup"
//
//	request := gorequest.New().Type("multipart").AppendHeader("Content-Type", "multipart/form-data")
//
//	var requestBody bytes.Buffer
//	multipartWriter := multipart.NewWriter(&requestBody)
//
//	var mediaData []map[string]interface{}
//
//	for index, file := range files {
//		src, err := file.Open()
//		if err != nil {
//			return "", err
//		}
//		defer src.Close()
//
//		request.SendFile(src, file.Filename, fmt.Sprintf("image%d", index))
//
//		mediaItem := map[string]interface{}{
//			"type":  "photo",
//			"media": fmt.Sprintf("attach://image%d", index),
//		}
//		mediaData = append(mediaData, mediaItem)
//	}
//
//
//	formData := map[string]string{
//		"chat_id": "@smm_auto_test",
//		"media": "[{\"type\":\"photo\", \"media\":\"attach://image0\"}, {\"media\":\"attach://image1\", \"type\":\"photo\"}]",
//	}
//	fmt.Println(formData)
//
//	for key, value := range formData {
//		_ = multipartWriter.WriteField(key, value)
//	}
//
//	multipartWriter.Close()
//	fmt.Println(requestBody.String())
//	_, body, errs := request.
//		SendString(requestBody.String()).
//		Post(url).
//		End()
//
//	if errs != nil {
//		fmt.Println("Error performing request:", errs)
//		return "", errs[0]
//	}
//	fmt.Println("Response body:", body)
//	return body, nil
//}

//func SendMediaGroup(files []*multipart.FileHeader, caption string) (message string, err error)  {
//	url := "https://api.telegram.org/bot" + os.Getenv("botToken") + "/sendMediaGroup"
//	client := resty.New()
//	req := client.R()
//
//	var formData = map[string]string{
//		"chat_id":"@smm_auto_test",
//		"media":"[{\"type\":\"photo\", \"media\":\"attach://image0\"}, {\"type\":\"photo\", \"media\":\"attach://image1\"}]",
//	}
//	request:= gorequest.New().Type("multipart")
//	for index, file := range files {
//		src, err := file.Open()
//		if err!=nil {
//			return "", err
//		}
//		var fileBuf bytes.Buffer
//		_, _ = io.Copy(&fileBuf, src)
//		req.SetFileReader(fmt.Sprintf("image%d", index), file.Filename, bytes.NewReader(fileBuf.Bytes()))
//		request.SendFile(fileBuf.Bytes(), file.Filename, fmt.Sprintf("image%d", index))
//	}
//	req.SetFormData(formData)
//	fmt.Println(req)
//	res, _, errrr := request.Send(formData).Post(url).End(func(response gorequest.Response, body string, errs []error){
//		if response.StatusCode != 200 {
//			fmt.Println("Error:", response.Status)
//		} else {
//			fmt.Println("Success:", response.Status)
//		}
//	})
//	if errrr!=nil {
//		return "", err
//	}
//
//	fmt.Println(res)
//
//	return "", nil
//}

//func SendMediaGroup(files []*multipart.FileHeader, caption string) (message string, err error)  {
//	url := "https://api.telegram.org/bot" + os.Getenv("botToken") + "/sendMediaGroup"
//	client:=resty.New()
//	var media = "["
//
//	body := &bytes.Buffer{}
//	writer := multipart.NewWriter(body)
//	defer writer.Close()
//
//	for index, file := range files {
//		fieldName := fmt.Sprintf("image%d", index)
//		imageField, err := writer.CreateFormFile(fieldName, file.Filename)
//		if err != nil {
//			return "", err
//		}
//		src, err := file.Open()
//		if err != nil {
//			return "", err
//		}
//		defer src.Close()
//		_, err = io.Copy(imageField, src)
//		if err != nil {
//			return "", err
//		}
//	}
//	for index, _ := range files {
//		newMedia := fmt.Sprintf(`{"type":"photo", "media":"attach://image%d"}, `,index)
//		media += newMedia
//		if index == len(files)-1 {
//			media = media[:len(media)-2]
//			media += "]"
//		}
//	}
//	_ = writer.WriteField("chat_id", "@smm_auto_test")
//	_ = writer.WriteField("media", media)
//	fmt.Println(string(body.Bytes()))
//	req:= client.R().SetBody(body.Bytes()).SetHeader("Content-Type", "multipart/form-data")
//
//	res, err := req.Post(url)
//	fmt.Println("response is ...")
//	fmt.Println(res.Status())
//	fmt.Println(res.String())
//	fmt.Println(string(res.Body()))
//
//	return "success", nil
//}

// SendPhoto TODO: Remake this shit into using io.Reader instead of *multipart.File
func SendPhoto(src *multipart.File, caption string, fileName string) (message string, err error) {
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
	_, err = io.Copy(imageField, *src)
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

