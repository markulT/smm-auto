package scrapper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

type UserAuthorizer interface {
	generateCode(cn string) Code
}

type defaultUserAuthorizer struct {
	chRepo ChannelAuthorizerRepository
}

func NewDefaultUserAuthorizer() UserAuthorizer {
	return &defaultUserAuthorizer{}
}


func (u *defaultUserAuthorizer) generateCode(cn string) Code {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(900000) + 100000
	return Code{
		ChannelName: cn,
		SubmitCode:  strconv.Itoa(code),
		ExpireDate:  time.Time{},
	}
}

type SendMessageRequest struct {
	Text                  string `json:"text"`
	ParseMode             string `json:"parse_mode"`
	DisableWebPagePreview bool   `json:"disable_web_page_preview"`
	DisableNotification   bool   `json:"disable_notification"`
	ReplyToMessage        string `json:"reply_to_message"`
	ChatId                string `json:"chat_id"`
}

func SendMessage(text string, chatId string) error {
	url := "https://api.telegram.org/bot" + os.Getenv("botToken") + "/sendMessage"
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
