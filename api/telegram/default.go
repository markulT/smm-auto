package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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