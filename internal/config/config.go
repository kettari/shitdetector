package config

import (
	"github.com/kettari/shitdetector/errors"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

type Config struct {
	BotToken string
	Debug    bool
}

func GetConfig() *Config {
	config := Config{}
	config.BotToken = os.Getenv("BOT_TOKEN")
	if len(config.BotToken) == 0 {
		log.Fatal(errors.ErrBotTokenIsEmpty)
	}

	debug := os.Getenv("BOT_DEBUG")
	if strings.ToLower(debug) == "true" {
		config.Debug = true
	}
	if config.Debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	return &config
}
