package main

import (
	"github.com/trydex/chatgpt-bot/internal/services"
)

func main() {

	chatGPTService := services.NewChatGPTService("YOUR_OPENAI_TOKEN")
	tgBot := services.NewTgBot("YOUR_BOT_TOKEN", chatGPTService)
	tgBot.Run()
}
