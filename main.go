package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	sendler "github.com/yudintsevegor/dotfiles/go_projects/src/tgBotVkPostsSendler"
	// sendler "github.com/yudintsevegor/tgBotVkPostsSendler"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func main() {

	db, err := sql.Open("postgres", DSN)
	if err != nil {
		log.Fatal("OPEN ERROR: ", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("PING ERROR: ", err)
	}

	// testing
	return

	port := "8080"
	go http.ListenAndServe(":"+port, nil)
	fmt.Printf("start listen :%v\n", port)

	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Fatal(err)
	}

	groupID := "-65652356" // https://vk.com/ustami_msu
	channelName := "@DebuggingMSUBot"
	webHookURL := "https://cfb24135.ngrok.io"

	opt := sendler.ReqOptions{
		Count:  "5",
		Offset: "0",
		Filter: "owner",
	}

	caller := sendler.Caller{
		ChannelName: channelName,
		WebHookURL:  webHookURL,
		Options:     opt,
		ErrChan:     make(chan error),
	}

	go caller.CallBot(bot, caller.GetVkPosts(groupID, ServiceKey))

	for err := range caller.ErrChan {
		log.Fatal(err)
	}
}
