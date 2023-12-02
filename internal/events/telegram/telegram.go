package telegram

import (
	"errors"
	"fmt"

	"github.com/AlexCorn999/link-saving-tgbot/internal/clients"
	"github.com/AlexCorn999/link-saving-tgbot/internal/clients/telegram"
	"github.com/AlexCorn999/link-saving-tgbot/internal/events"
	"github.com/AlexCorn999/link-saving-tgbot/internal/storage"
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID   int
	UserName string
}

func NewProcessor(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, fmt.Errorf("can't get events: %w", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return errors.New("unknown event type")
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return fmt.Errorf("can't process message : %w", err)
	}

	if err := p.doCmd(meta.ChatID, meta.UserName, event.Text); err != nil {
		return fmt.Errorf("can't process message : %w", err)
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, errors.New("can't get meta. Unknown meta type")
	}
	return res, nil
}

func event(upd clients.Update) events.Event {
	updType := fetchType(upd)
	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			UserName: upd.Message.From.UserName,
		}
	}

	return res
}

func fetchText(upd clients.Update) string {
	if upd.Message == nil {
		return ""
	}
	return upd.Message.Text
}

func fetchType(upd clients.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}
	return events.Message
}
