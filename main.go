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
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("PING ERROR: ", err)
	}

	w := sendler.Writer{
		DB: db,
	}

	tableName := "quotesMSU"
	if _, err = w.EditTable(sendler.CreateTable, tableName); err != nil {
		log.Fatal("EDITION ERROR: ", err)
	}

	if err = w.WriteToDB(tableName); err != nil {
		log.Fatal("WRITING ERROR: ", err)
	}

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
