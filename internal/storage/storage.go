package storage

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
)

var ErrorNoSavedPages = errors.New("no saved pages")

type Storage interface {
	Save(ctx context.Context, p *Page) error
	PickRandom(ctx context.Context, userName string) (*Page, error)
	Remove(ctx context.Context, p *Page) error
	IsExists(ctx context.Context, p *Page) (bool, error)
}

type Page struct {
	URL      string
	UserName string
}

func (p *Page) Hash() (string, error) {
	hash := sha1.New()

	if _, err := hash.Write([]byte(p.URL)); err != nil {
		return "", fmt.Errorf("can't calculate hash: %w", err)
	}

	return fmt.Sprintf("%x", hash.Sum([]byte(p.UserName))), nil
}
