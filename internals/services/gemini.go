package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
)

var APIKeyNotFound = errors.New("api key not found for gemini")

type GeminiService struct {
	apiKey string
	client *genai.Client
}

func NewGeminiService() (*GeminiService, error) {
	apiKey := viper.GetString("ai.gemini_api_key")
	ctx := context.Background()
	if apiKey == "" {
		return nil, APIKeyNotFound
	}
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	return &GeminiService{
		apiKey: apiKey,
		client: client,
	}, nil
}

func (gs GeminiService) GenerateText(prompt string) ([]string, error) {
	model := gs.client.GenerativeModel("gemini-pro")
	resp, err := model.GenerateContent(context.Background(), genai.Text("prompt"))
	if err != nil {
		return []string{}, err
	} else {
		responses := []string{}
		for _, candidate := range resp.Candidates {
			response := ""
			for _, part := range candidate.Content.Parts {
				response = fmt.Sprintf("%+v %+v", response, part)
			}
			responses = append(responses, response)
		}
		return responses, nil
	}
}

func (gs GeminiService) CloseClient() error {
	return gs.client.Close()
}
