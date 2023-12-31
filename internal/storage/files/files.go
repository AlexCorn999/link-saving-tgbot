package files

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/AlexCorn999/link-saving-tgbot/internal/storage"
)

const defaultPerm = 0774

type Storage struct {
	basePath string
}

func NewStorage(basePath string) *Storage {
	return &Storage{
		basePath: basePath,
	}
}

// Save creates a directory for the user and writes the link to a separate file in the user's directory.
func (s *Storage) Save(ctx context.Context, page *storage.Page) error {
	fPath := filepath.Join(s.basePath, page.UserName)

	if err := os.MkdirAll(fPath, defaultPerm); err != nil {
		return fmt.Errorf("can't save page: %w", err)
	}

	fName, err := fileName(page)
	if err != nil {
		return fmt.Errorf("can't save page: %w", err)
	}

	fPath = filepath.Join(fPath, fName)

	file, err := os.Create(fPath)
	if err != nil {
		return fmt.Errorf("can't save page: %w", err)
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(page); err != nil {
		return fmt.Errorf("can't save page: %w", err)
	}

	return nil
}

// PickRandom takes a random link from the user's storage.
func (s *Storage) PickRandom(ctx context.Context, userName string) (*storage.Page, error) {
	fPath := filepath.Join(s.basePath, userName)

	switch _, err := os.Stat(fPath); {
	case errors.Is(err, os.ErrNotExist):
		if err := os.MkdirAll(fPath, defaultPerm); err != nil {
			return nil, fmt.Errorf("can't make directory: %w", err)
		}
	case err != nil:
		return nil, fmt.Errorf("can't check directory %s", fPath)
	}

	files, err := os.ReadDir(fPath)
	if err != nil {
		return nil, fmt.Errorf("can't pick random page: %w", err)
	}

	if len(files) == 0 {
		return nil, storage.ErrorNoSavedPages
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := r.Intn(len(files))
	file := files[n]
	return s.decodePage(filepath.Join(fPath, file.Name()))
}

// Remove deletes the file from the user's storage.
func (s *Storage) Remove(ctx context.Context, p *storage.Page) error {
	fileName, err := fileName(p)
	if err != nil {
		return fmt.Errorf("can't remove file: %w", err)
	}

	fPath := filepath.Join(s.basePath, p.UserName, fileName)
	if err := os.Remove(fPath); err != nil {
		return fmt.Errorf("can't remove file %s : %w", fPath, err)
	}
	return nil
}

// IsExists checks if a file exists in the user's storage.
func (s *Storage) IsExists(ctx context.Context, p *storage.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, fmt.Errorf("can't check if file exists: %w", err)
	}

	fPath := filepath.Join(s.basePath, p.UserName, fileName)

	switch _, err := os.Stat(fPath); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("can't check if file %s exists", fPath)
	}

	return true, nil
}

// decodePage translates the file data from JSON to a structure.
func (s *Storage) decodePage(filePath string) (*storage.Page, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("can't decode page: %w", err)
	}
	defer file.Close()

	var p storage.Page

	if err := json.NewDecoder(file).Decode(&p); err != nil {
		return nil, fmt.Errorf("can't decode page: %w", err)
	}
	return &p, nil
}

// fileName gets the name of the hashed file.
func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
