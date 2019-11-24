package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	sendler "github.com/yudintsevegor/dotfiles/go_projects/src/tgBotVkPostsSendler"
	// sendler "github.com/yudintsevegor/tgBotVkPostsSendler"
)

func main() {
	db, err := sql.Open("postgres", DSN)
	if err != nil {
		log.Fatal("DataBase Open error: ", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("DataBase Ping Error: ", err)
	}

	w := sendler.DbWriter{
		DB:        db,
		TableName: "quotesmsu",
	}

	if _, err = w.CreateTable(); err != nil {
		log.Fatal("Createtable Error: ", err)
	}

	port := "8080"
	go http.ListenAndServe(":"+port, nil)
	fmt.Printf("start listen :%v\n", port)

	groupID := "-65652356" // https://vk.com/ustami_msu
	channelName := "@DebuggingMSUBot"
	webHookURL := "https://ecd153a5.ngrok.io"

	telegram := sendler.Telegram{
		ChannelName: channelName,
		WebHookURL:  webHookURL,
		BotToken:    BotToken,
	}

	opt := sendler.ReqOptions{
		Count:  "10",
		Offset: "0",
		Filter: "owner",
	}

	handler := sendler.Handler{
		Telegram: telegram,
		Options:  opt,
		ErrChan:  make(chan error),

		TimeOut:  time.Hour * 24,
		DbWriter: &w,
	}

	recipients := []string{"georgesyndicart"}
	handler.GetRecipients(recipients)

	go handler.StartBot(handler.GetVkPosts(groupID, VkServiceKey))

	for err := range handler.ErrChan {
		log.Fatal(err)
	}
}
