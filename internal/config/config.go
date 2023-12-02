package config

import (
	"errors"
	"flag"
)

type Config struct {
	Token       string
	StoragePath string
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) MustToken() error {
	token := flag.String("tgToken",
		"",
		"token for access to telegram bot")
	storagePath := flag.String("storage-path",
		"files_storage",
		"file path for storage")
	flag.Parse()

	c.Token = *token
	c.StoragePath = *storagePath

	if c.Token == "" {
		return errors.New("token is not specified")
	}
	return nil
}
