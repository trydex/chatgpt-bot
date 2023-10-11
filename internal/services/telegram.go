package services

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	NeedAnswerPrefix                 = "брат"
	NeedAnswerWithClearContextPrefix = "братан"
)

type tgBot struct {
	bot            *tgbotapi.BotAPI
	chatGptService ChatGPTService
}

type TgBot interface {
	Run()
}

func NewTgBot(token string, chatGptService ChatGPTService) TgBot {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	return &tgBot{
		bot:            bot,
		chatGptService: chatGptService,
	}

}

func (tgBot *tgBot) Run() {

	fmt.Println("Bot running")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := tgBot.bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil && needAnswer(update.Message.Text) {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			_, _ = tgBot.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Брат, дай минутку подумать над ответом"))

			fmt.Printf("Send request to ChatGPT \n")
			resp, err := tgBot.chatGptService.CreateChatCompletion(update.Message.Chat.ID, update.Message.Text, needClearContext(update.Message.Text))

			fmt.Printf("[ChatGPT] %s", resp.Choices[0].Message.Content)
			if err != nil {
				fmt.Printf("ChatCompletion error: %v\n", err)
				_, _ = tgBot.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Брат, произошла ошибка. Повтори запрос."))
				continue
			}

			fmt.Printf("[ChatGPT] %s", resp.Choices[0].Message.Content)

			replied := false
			for _, msgText := range splitLargeMessage(resp.Choices[0].Message.Content) {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
				if !replied {
					msg.ReplyToMessageID = update.Message.MessageID
					replied = true
				}

				_, _ = tgBot.bot.Send(msg)
			}

		}
	}

	fmt.Println("Bot stopped")
}

func needAnswer(msg string) bool {
	lowerMsg := strings.ToLower(msg)
	return strings.HasPrefix(lowerMsg, NeedAnswerPrefix) || strings.HasPrefix(lowerMsg, NeedAnswerWithClearContextPrefix)
}

func needClearContext(msg string) bool {
	lowerMsg := strings.ToLower(msg)
	return strings.HasPrefix(lowerMsg, NeedAnswerWithClearContextPrefix)
}

func splitLargeMessage(msg string) []string {
	return chunkString(msg, 2000)
}

func chunkString(s string, chunkSize int) []string {
	var chunks []string
	runes := []rune(s)

	if len(runes) == 0 {
		return []string{s}
	}

	for i := 0; i < len(runes); i += chunkSize {
		nn := i + chunkSize
		if nn > len(runes) {
			nn = len(runes)
		}
		chunks = append(chunks, string(runes[i:nn]))
	}
	return chunks
}
