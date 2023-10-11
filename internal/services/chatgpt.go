package services

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

type chatGptService struct {
	client          *openai.Client
	historyMessages map[int64][]openai.ChatCompletionMessage
}

type ChatGPTService interface {
	CreateChatCompletion(chatId int64, content string, clearContext bool) (openai.ChatCompletionResponse, error)
}

func NewChatGPTService(token string) ChatGPTService {
	return &chatGptService{
		client:          openai.NewClient(token),
		historyMessages: make(map[int64][]openai.ChatCompletionMessage),
	}
}

func (s *chatGptService) CreateChatCompletion(chatId int64, content string, clearContext bool) (resp openai.ChatCompletionResponse, err error) {
	_, ok := s.historyMessages[chatId]
	if !ok || clearContext {
		s.historyMessages[chatId] = make([]openai.ChatCompletionMessage, 0)
	}

	s.historyMessages[chatId] = append(s.historyMessages[chatId], openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: content,
	})

	resp, err = s.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: s.historyMessages[chatId],
		},
	)

	if err == nil {
		s.historyMessages[chatId] = append(s.historyMessages[chatId], openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: resp.Choices[0].Message.Content,
		})
	}

	return resp, err
}
