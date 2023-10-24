package main

import (
	"github.com/trydex/chatgpt-bot/internal/config"
	"github.com/trydex/chatgpt-bot/internal/services"
	"log"
)

func main() {

	if err := config.Load(); err != nil {
		log.Fatalf("config wasn't loaded: %v", err)
	}

	cfg, err := config.NewAppConfig()
	if err != nil {
		log.Fatalf("app config wasn't loaded: %v", err)
	}

	chatGPTService := services.NewChatGPTService(cfg.OpenAIKey())
	tgBot := services.NewTgBot(cfg.TgBotKey(), chatGPTService)
	tgBot.Run()
}
