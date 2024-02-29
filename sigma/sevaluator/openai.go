package sevaluator

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

// OpenAIService interface defines the operations for OpenAI service
type OpenAIService interface {
	CreateChatCompletion(context.Context, openai.ChatCompletionRequest) (*openai.ChatCompletionResponse, error)
}

// OpenAIClient implements OpenAIService
type OpenAIClient struct {
	*openai.Client
}

// NewOpenAIClient creates a new OpenAIClient instance
func NewOpenAIClient(apiKey string) *OpenAIClient {
	return &OpenAIClient{
		Client: openai.NewClient(apiKey),
	}
}

// SendMessageToOpenAI sends a message to OpenAI and returns the response
func SendMessageToOpenAI(apiKey, content string) (string, error) {
	openAIClient := NewOpenAIClient(apiKey)

	request := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: content,
			},
		},
	}

	resp, err := openAIClient.CreateChatCompletion(context.Background(), request)
	if err != nil {
		return "", fmt.Errorf("ChatCompletion error: %v", err)
	}

	if len(resp.Choices) > 0 {
		return resp.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("empty response from OpenAI")
}
