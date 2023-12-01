package config

import (
	"errors"
	"flag"
)

type Config struct {
	Token string
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) MustToken() error {
	c.Token = *flag.String("token-bot-token",
		"",
		"token for access to telegram bot")

	flag.Parse()

	if c.Token == "" {
		return errors.New("token is not specified")
	}
	return nil
}
