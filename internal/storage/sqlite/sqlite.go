package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"github.com/AlexCorn999/link-saving-tgbot/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func (s *Storage) Init(ctx context.Context) error {
	if _, err := s.db.ExecContext(ctx,
		"CREATE TABLE IF NOT EXISTS pages (url TEXT, user_name TEXT)"); err != nil {
		return fmt.Errorf("can't create table: %w", err)
	}
	return nil
}

func NewStorage(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}

	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) Save(ctx context.Context, p *storage.Page) error {
	if _, err := s.db.ExecContext(ctx, "INSERT INTO pages (url, user_name) VALUES (?, ?)",
		p.URL,
		p.UserName); err != nil {
		return fmt.Errorf("can't save page: %w", err)
	}
	return nil
}

func (s *Storage) PickRandom(ctx context.Context, userName string) (*storage.Page, error) {
	var url string
	if err := s.db.QueryRowContext(ctx, "SELECT url FROM pages WHERE user_name = ? ORDER BY RANDOM() LIMIT 1",
		userName).Scan(&url); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrorNoSavedPages
		} else {
			return nil, fmt.Errorf("can't pick random page: %w", err)
		}
	}

	return &storage.Page{
		URL:      url,
		UserName: userName,
	}, nil
}

func (s *Storage) Remove(ctx context.Context, p *storage.Page) error {
	if _, err := s.db.ExecContext(ctx, "DELETE FROM pages WHERE url = ? AND user_name = ?",
		p.URL,
		p.UserName); err != nil {
		return fmt.Errorf("can't remove page: %w", err)
	}
	return nil
}
func (s *Storage) IsExists(ctx context.Context, p *storage.Page) (bool, error) {
	var count int
	if err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM pages WHERE url = ? AND user_name = ?",
		p.URL, p.UserName).Scan(&count); err != nil {
		return false, fmt.Errorf("can't check if page exists: %w", err)
	}

	return count > 0, nil
}
