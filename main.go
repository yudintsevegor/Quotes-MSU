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
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

const (
	isExist = `
	SELECT EXISTS (
		SELECT 1
		FROM   information_schema.tables 
		WHERE  table_schema = 'public'
		AND    table_name = 'quotesmsu'
		);`
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
		DB:        db,
		TableName: "quotesMSU",
	}

	var ok bool
	row := w.DB.QueryRow(isExist)
	row.Scan(&ok)
	log.Println(ok)

	/**/
	if !ok{
		if _, err = w.EditTable(sendler.CreateTable); err != nil {
			log.Fatal("EDITION DB ERROR: ", err)
		}
	}
	
	/**/

	port := "8080"
	go http.ListenAndServe(":"+port, nil)
	fmt.Printf("start listen :%v\n", port)

	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Fatal(err)
	}

	groupID := "-65652356" // https://vk.com/ustami_msu
	channelName := "@DebuggingMSUBot"
	webHookURL := "https://9a37a229.ngrok.io"

	opt := sendler.ReqOptions{
		Count:  "10",
		Offset: "0",
		Filter: "owner",
	}

	caller := sendler.Caller{
		ChannelName: channelName,
		WebHookURL:  webHookURL,
		Options:     opt,
		ErrChan:     make(chan error),

		TimeOut: time.Hour * 24,
		Writer:  &w,
	}

	go caller.CallBot(bot, caller.GetVkPosts(groupID, ServiceKey))

	for err := range caller.ErrChan {
		log.Fatal(err)
	}
}
