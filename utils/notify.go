package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"rss-reader/globals"
	"strings"
)

const (
	FeiShuRoute   = "feishu"
	TelegramRoute = "telegram"
	ContentType   = "application/json"
	TokenReplace  = "${token}"
)

type Message struct {
	Routes  []string `json:"routes"`
	Content string   `json:"content"`
}

type FeiShuMessage struct {
	MsgType string            `json:"msg_type"`
	Content FeiShuMessageText `json:"content"`
}

type FeiShuMessageText struct {
	Text string `json:"text"`
}

type TelegramMessage struct {
	ChatId string `json:"chat_id"`
	Text   string `json:"text"`
}

func Notify(msg Message) {
	if msg.Routes == nil || len(msg.Routes) == 0 {
		return
	}
	for _, route := range msg.Routes {
		switch route {
		case FeiShuRoute:
			if globals.RssUrls.Notify.FeiShu.API != "" {
				sendToFeiShu(msg)
			}
		case TelegramRoute:
			if globals.RssUrls.Notify.Telegram.API != "" {
				sendToTelegram(msg)
			}
		default:
			log.Println("without route")
		}
	}
}
func sendToTelegram(msg Message) {
	finalMsg, err := json.Marshal(
		TelegramMessage{
			ChatId: globals.RssUrls.Notify.Telegram.ChatId,
			Text:   msg.Content,
		})
	if err != nil {
		log.Printf("json marshal err: %+v\n", err)
		return
	}
	api := strings.ReplaceAll(globals.RssUrls.Notify.Telegram.API, TokenReplace, globals.RssUrls.Notify.Telegram.Token)
	requestPost(api, finalMsg)
}

func sendToFeiShu(msg Message) {
	finalMsg, err := json.Marshal(
		FeiShuMessage{
			MsgType: "text",
			Content: FeiShuMessageText{
				Text: msg.Content,
			},
		})
	if err != nil {
		log.Printf("json marshal err: %+v\n", err)
		return
	}
	requestPost(globals.RssUrls.Notify.FeiShu.API, finalMsg)
}

func requestPost(url string, param []byte) {
	requestBody := bytes.NewBuffer(param)
	resp, err := http.Post(url, ContentType, requestBody)

	if err != nil {
		log.Printf("http post err: %+v\n", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("http body close err: %+v\n", err)
		}
	}(resp.Body)
	body, err := ioutil.ReadAll(resp.Body) // 读取响应内容
	if err != nil {
		log.Printf("http post read body err: %+v\n", err)
		return
	}
	log.Printf("response status: %s,response body:%s", string(body), resp.Status)
	//string(body)
	return
}
