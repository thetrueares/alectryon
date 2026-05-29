package provider

import (
	"context"

	"go.iain.rocks/alectryon/backend/engine"
	"google.golang.org/genai"
)

func NewGemini(ctx context.Context, apiKey string) (*GeminiAi, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})

	if err != nil {
		return nil, err
	}

	return &GeminiAi{client: client}, nil
}

type GeminiAi struct {
	client *genai.Client
}

func (ga GeminiAi) Process(in engine.Input) engine.Output {

	response, err := ga.client.Models.GenerateContent(context.TODO(), "gemini-2.0-flash", genai.Text(in.Text), nil)

	if err != nil {
		return engine.Output{
			Text:       err.Error(),
			TokenCount: 1,
		}
	}

	return engine.Output{
		Text:       response.Text(),
		TokenCount: 10,
	}
}
