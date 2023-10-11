package main

import (
	"github.com/trydex/chatgpt-bot/internal/services"
)

func main() {

	chatGPTService := services.NewChatGPTService("YOUR_KEY")
	tgBot := services.NewTgBot("YOUR_KEY", chatGPTService)
	tgBot.Run()
}
