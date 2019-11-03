package main

/*
import (
	"fmt"
	"net/http"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)


func callBot(bot *tgbotapi.BotAPI, in <-chan string, channelName string) chan error {
	errChan := make(chan error)
	// bot.Debug = true
	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)

	_, err := bot.SetWebhook(tgbotapi.NewWebhook(WebHookURL))
	if err != nil {
		errChan <- err
	}

	updates := bot.ListenForWebhook("/")

	port := "8080"
	go http.ListenAndServe(":"+port, nil)
	fmt.Printf("start listen :%v", port)

	go func() {
		for {
			select {
			case text := <-in:
				if _, err := bot.Send(tgbotapi.NewMessageToChannel(channelName, text)); err != nil {
					errChan <- err
				}
			case update := <-updates:
				_, err = bot.Send(tgbotapi.NewMessage(
					update.Message.Chat.ID,
					fmt.Sprintf("Bot is handler for %v channel", channelName),
				))
				if err != nil {
					errChan <- err
				}
			}
		}
	}()

	return errChan
}
*/
