package telegram

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/AlexCorn999/link-saving-tgbot/internal/storage"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

// doCmd processes commands.
func (p *Processor) doCmd(chatID int, username, text string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command %s from %s", text, username)

	if isAddCmd(text) {
		return p.savePage(chatID, username, text)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

// savePage checks if the link already exists in the repository and if not, saves it.
func (p *Processor) savePage(chatID int, username, pageURL string) error {
	page := storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	isExists, err := p.storage.IsExists(context.Background(), &page)
	if err != nil {
		return err
	}

	if isExists {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}

	if err := p.storage.Save(context.Background(), &page); err != nil {
		return err
	}

	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}

	return nil
}

// sendRandom sends a random link from the user's storage and deletes it.
func (p *Processor) sendRandom(chatID int, username string) error {

	page, err := p.storage.PickRandom(context.Background(), username)
	if err != nil && !errors.Is(err, storage.ErrorNoSavedPages) {
		return fmt.Errorf("can't send random page: %w", err)
	}

	if errors.Is(err, storage.ErrorNoSavedPages) {
		return fmt.Errorf("can't send random page: %w", p.tg.SendMessage(chatID, msgNoSavedPages))
	}

	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return fmt.Errorf("can't send random page: %w", err)
	}

	return p.storage.Remove(context.Background(), page)
}

// sendHelp sends information to help.
func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

// sendHello sends a welcome message.
func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

// isAddCmd checks for adding a link.
func isAddCmd(text string) bool {
	return isURL(text)
}

// isURL checks to see if the url.
func isURL(text string) bool {
	u, err := url.Parse(text)
	return err == nil && u.Host != ""
}
