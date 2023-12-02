package main

import (
	"log"

	"github.com/AlexCorn999/link-saving-tgbot/internal/clients/telegram"
	"github.com/AlexCorn999/link-saving-tgbot/internal/config"
)

const (
	tgBotHost = "api.telegram.org"
)

func main() {
	config := config.NewConfig()
	if err := config.MustToken(); err != nil {
		log.Fatal(err)
	}

	tgClient := telegram.NewClient(tgBotHost, config)
	//fetcher := fetcher.New()
	//processor := telegram.NewProcessor()
	//consumer.Start(fetcher,processor)
}
