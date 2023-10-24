package config

import (
	"errors"
	"flag"
	"github.com/joho/godotenv"
	"os"
)

const (
	tgBotKey  = "TGBOT_KEY"
	openAIKey = "OPENAI_KEY"
)

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

var configPath string

type AppConfig interface {
	TgBotKey() string
	OpenAIKey() string
}

type appConfig struct {
	tgBotKey  string
	openAIKey string
}

func Load() error {
	flag.Parse()

	err := godotenv.Load(configPath)
	if err != nil {
		return err
	}

	return nil
}

func NewAppConfig() (AppConfig, error) {
	tgBotKey := os.Getenv(tgBotKey)
	if len(tgBotKey) == 0 {
		return nil, errors.New("tg bot key not found")
	}

	openAIKey := os.Getenv(openAIKey)
	if len(openAIKey) == 0 {
		return nil, errors.New("open ai key not found")
	}

	return &appConfig{
		openAIKey: openAIKey,
		tgBotKey:  tgBotKey,
	}, nil
}

func (a *appConfig) TgBotKey() string {
	return a.tgBotKey
}

func (a *appConfig) OpenAIKey() string {
	return a.openAIKey
}
