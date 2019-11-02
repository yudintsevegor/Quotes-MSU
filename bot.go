package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

const (
	WebHookURL = "https://94d76e27.ngrok.io"
)

var rss = map[string]string{
	"Habr": "https://habrahabr.ru/rss/best/",
}

type RSS struct {
	Items []Item `xml:"channel>item"`
}

type Item struct {
	URL   string `xml:"guid"`
	Title string `xml:"title"`
}

func getNews(url string) (*RSS, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	rss := new(RSS)
	if err = xml.Unmarshal(body, rss); err != nil {
		return nil, err
	}

	return rss, nil
}

const text = "Quotes"

func bot(path string) {
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Fatal(err)
	}

	// bot.Debug = true
	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(WebHookURL))
	if err != nil {
		log.Fatal(err)
	}

	updates := bot.ListenForWebhook("/")

	// port := os.Getenv("PORT")
	port := "8080"
	go http.ListenAndServe(":"+port, nil)
	fmt.Println("start listen :8080")

	// получаем все обновления из канала updates
	for update := range updates {
		log.Println("HERE")
		/**
		url, ok := rss[update.Message.Text]
		if !ok {
			bot.Send(tgbotapi.NewMessage(
				update.Message.Chat.ID,
				`there is only 'Habr' feed availible`,
			))
			continue
		}
		/**
		rss, err := getNews(url)
		if err != nil {
			log.Println(err)
			bot.Send(tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"sorry, error happend",
			))
			continue
		}
			/**/
		if update.Message.Text != text {
			bot.Send(tgbotapi.NewMessage(
				update.Message.Chat.ID,
				fmt.Sprintf("there is only %v feed availible", text),
			))
			continue
		}

		body, err := getPosts(path)
		if err != nil {
			log.Println(err)
			bot.Send(tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"sorry, error happend",
			))
			continue
		}

		if len(body.Groups) != 1 {
			errText := "empty info about group"
			err := errors.New(errText)
			log.Println(err)
			bot.Send(tgbotapi.NewMessage(
				update.Message.Chat.ID,
				errText,
			))
			continue
		}

		for _, v := range body.Items {
			bot.Send(tgbotapi.NewMessage(
				update.Message.Chat.ID,
				v.Text+"\n"+makeLink(strconv.Itoa(v.ID)),
			))
			// log.Println(time.Time(v.Date))
		}
		/**/
		/**
		for _, item := range rss.Items {
			bot.Send(tgbotapi.NewMessage(
				update.Message.Chat.ID,
				item.URL+"\n"+item.Title,
			))
		}
		/**/

	}
}
