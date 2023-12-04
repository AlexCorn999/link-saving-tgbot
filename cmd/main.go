package main

import (
	"context"
	"log"

	tgClient "github.com/AlexCorn999/link-saving-tgbot/internal/clients/telegram"
	"github.com/AlexCorn999/link-saving-tgbot/internal/config"
	event_consumer "github.com/AlexCorn999/link-saving-tgbot/internal/consumer/event-consumer"
	"github.com/AlexCorn999/link-saving-tgbot/internal/events/telegram"
	"github.com/AlexCorn999/link-saving-tgbot/internal/storage/sqlite"
)

const (
	tgBotHost = "api.telegram.org"
)

func main() {
	config := config.NewConfig()
	if err := config.ParseFlags(); err != nil {
		log.Fatal(err)
	}

	//store := files.NewStorage(config.FilePath)

	store, err := sqlite.NewStorage(config.SqlitePath)
	if err != nil {
		log.Fatal("can't connect to storage:", err)
	}

	if err := store.Init(context.Background()); err != nil {
		log.Fatal("can't init storage")
	}

	eventsProcessor := telegram.NewProcessor(tgClient.NewClient(tgBotHost, config),
		store)

	log.Println("server started")

	consumer := event_consumer.NewConsumer(eventsProcessor, eventsProcessor, 100)
	if err := consumer.Start(); err != nil {
		log.Fatal(err)
	}
}
