package services

import (
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	NeedAnswerPrefix                 = "брат"
	NeedAnswerWithClearContextPrefix = "братан"
	UsGPT4Prefix                     = "братуха"
	TgMsgMaxLen                      = 2000
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

			go answerUserQuestion(tgBot, update.Message)
		}
	}

	fmt.Println("Bot stopped")
}

func answerUserQuestion(tgBot *tgBot, msg *tgbotapi.Message) {
	defer func() {
		if err := recover(); err != nil {
			askRepeatRequest(tgBot, msg.Chat.ID, "Братан, повтори запрос.")
		}
	}()

	confirmationMsg := tgbotapi.NewMessage(msg.Chat.ID, chooseConfirmationMsgText())
	confirmationMsg.ReplyToMessageID = msg.MessageID

	_, _ = tgBot.bot.Send(confirmationMsg)

	fmt.Printf("Send request to ChatGPT \n")
	resp, err := tgBot.chatGptService.CreateChatCompletion(msg.Chat.ID, msg.Text, needClearContext(msg.Text), useGPT4(msg.Text))
	if err != nil {
		askRepeatRequest(tgBot, msg.Chat.ID, err.Error())
		return
	}

	fmt.Printf("[ChatGPT] %s", resp.Choices[0].Message.Content)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		askRepeatRequest(tgBot, msg.Chat.ID, err.Error())
		return
	}

	fmt.Printf("[ChatGPT] %s", resp.Choices[0].Message.Content)

	replied := false
	for _, msgText := range splitLargeMessage(resp.Choices[0].Message.Content) {
		message := tgbotapi.NewMessage(msg.Chat.ID, replaceSpecialCharacters(msgText))
		message.ParseMode = "MarkdownV2"

		if !replied {
			message.ReplyToMessageID = msg.MessageID
			replied = true
		}

		_, err = tgBot.bot.Send(message)
		if err != nil {
			_, err = tgBot.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, msgText))

			if err != nil {
				askRepeatRequest(tgBot, msg.Chat.ID, err.Error())
			}
		}
	}
}

func askRepeatRequest(tgBot *tgBot, chatId int64, text string) {
	_, _ = tgBot.bot.Send(tgbotapi.NewMessage(chatId, text))
}

func needAnswer(msg string) bool {
	lowerMsg := strings.ToLower(msg)
	return strings.HasPrefix(lowerMsg, NeedAnswerPrefix) || strings.HasPrefix(lowerMsg, NeedAnswerWithClearContextPrefix)
}

func useGPT4(msg string) bool {
	lowerMsg := strings.ToLower(msg)
	return strings.HasPrefix(lowerMsg, UsGPT4Prefix)
}

func replaceSpecialCharacters(msg string) string {
	re := regexp.MustCompile(`([|{\[\]*_~}+)(#>!=\-.])`)

	return re.ReplaceAllString(msg, "\\$1")
}

var confirmationMessages = []string{
	"Брат, дай минутку подумать над ответом",
	"Братан, подожди секунду, пережевываю вопрос",
	"Братуха, дай мне немного времени на раздумье",
	"Братишка, давай задание переварю, точно найду ответ",
	"Браток, дай мне минутку, чтобы собраться с мыслями",
	"Бро, дай мне пару секунд, чтобы взвесить все варианты ответов",
	"Братанчик, мне нужно пару мгновений, чтобы подумать над ответом",
	"Брателло, удели мне несколько моментов для размышления",
	"Братишка, полагаю, я смогу ответить после короткой паузы",
	"Брат, дай мне секунду, чтобы сообразить, как отвечать",
}

func chooseConfirmationMsgText() string {
	return confirmationMessages[rand.Intn(len(confirmationMessages))]
}

func needClearContext(msg string) bool {
	lowerMsg := strings.ToLower(msg)
	return strings.HasPrefix(lowerMsg, NeedAnswerWithClearContextPrefix)
}

func splitLargeMessage(msg string) []string {
	return chunkString(msg, TgMsgMaxLen)
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
