package config

import (
	"errors"
	"flag"
)

type Config struct {
	TgToken    string
	FilePath   string
	SqlitePath string
}

func NewConfig() *Config {
	return &Config{}
}

// ParseFlags parses the token for Telegram, the path for the filesystem, and the path for the sqlite database.
func (c *Config) ParseFlags() error {
	token := flag.String("tgToken",
		"",
		"token for access to telegram bot")
	fileStoragePath := flag.String("file",
		"data/files",
		"path for file storage")
	sqlitePath := flag.String("sqlite",
		"data/sqlite/storage.db",
		"path for sqlite storage")
	flag.Parse()

	c.TgToken = *token
	c.FilePath = *fileStoragePath
	c.SqlitePath = *sqlitePath

	if c.TgToken == "" {
		return errors.New("token is not specified")
	}
	return nil
}
