package main

import (
	"log"

	tgClient "github.com/AlexCorn999/link-saving-tgbot/internal/clients/telegram"
	"github.com/AlexCorn999/link-saving-tgbot/internal/config"
	event_consumer "github.com/AlexCorn999/link-saving-tgbot/internal/consumer/event-consumer"
	"github.com/AlexCorn999/link-saving-tgbot/internal/events/telegram"

	"github.com/AlexCorn999/link-saving-tgbot/internal/storage/files"
)

const (
	tgBotHost = "api.telegram.org"
)

func main() {
	config := config.NewConfig()
	if err := config.MustToken(); err != nil {
		log.Fatal(err)
	}

	eventsProcessor := telegram.NewProcessor(tgClient.NewClient(tgBotHost, config),
		files.NewStorage(config.StoragePath))

	log.Println("server started")

	consumer := event_consumer.NewConsumer(eventsProcessor, eventsProcessor, 100)
	if err := consumer.Start(); err != nil {
		log.Fatal(err)
	}
}
