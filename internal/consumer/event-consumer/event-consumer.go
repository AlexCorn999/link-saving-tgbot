package event_consumer

import (
	"log"
	"time"

	"github.com/AlexCorn999/link-saving-tgbot/internal/events"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func NewConsumer(fetcher events.Fetcher, processor events.Processor, batchSize int) *Consumer {
	return &Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

// Start starts the bot process.
func (c *Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())
			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(time.Second * 1)
			continue
		}

		if err := c.handleEvents(gotEvents); err != nil {
			log.Println(err)
			continue
		}
	}
}

// handleEvents passes events to processing.
func (c *Consumer) handleEvents(events []events.Event) error {
	for _, e := range events {
		log.Printf("got new event: %s", e.Text)

		if err := c.processor.Process(e); err != nil {
			log.Printf("can't handle event %s", err.Error())
			continue
		}
	}
	return nil
}
