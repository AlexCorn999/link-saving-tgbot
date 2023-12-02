package config

import (
	"errors"
	"flag"
)

type Config struct {
	Token           string
	FileStoragePath string
	SqlitePath      string
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) MustToken() error {
	token := flag.String("tgToken",
		"",
		"token for access to telegram bot")
	fileStoragePath := flag.String("storage-path",
		"files_storage",
		"file path for storage")
	sqlitePath := flag.String("sqlite",
		"data/sqlite/storage.db",
		"path for storage")
	flag.Parse()

	c.Token = *token
	c.FileStoragePath = *fileStoragePath
	c.SqlitePath = *sqlitePath

	if c.Token == "" {
		return errors.New("token is not specified")
	}
	return nil
}
