package storage

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
)

var ErrorNoSavedPages = errors.New("no saved pages")

type Storage interface {
	Save(p *Page) error
	PickRandom(userName string) (*Page, error)
	Remove(p *Page) error
	IsExists(p *Page) (bool, error)
}

type Page struct {
	URL      string
	UserName string
}

func (p *Page) Hash() (string, error) {
	h := sha256.New()
	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", fmt.Errorf("can't calculate hash: %w", err)
	}

	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", fmt.Errorf("can't calculate hash: %w", err)
	}

	return string(h.Sum(nil)), nil
}
