package vendor

import (
	"context"

	ollama "github.com/ollama/ollama/api"
	"go.iain.rocks/alectryon/api/engine"
)

func NewOllama(apiKey string) *Ollama {
	client, _ := ollama.ClientFromEnvironment()

	return &Ollama{client: client}
}

type Ollama struct {
	client *ollama.Client
}

func (oa Ollama) Process(input engine.Input) engine.Output {
	stream := false
	req := &ollama.GenerateRequest{
		Model:  "gemma4",
		Prompt: input.Text,
		Stream: &stream,
	}

	var output engine.Output

	err := oa.client.Generate(context.TODO(), req, func(resp ollama.GenerateResponse) error {
		output.Text += resp.Response
		return nil
	})

	if err != nil {
		return engine.Output{
			Text:       "Ollama error: " + err.Error(),
			TokenCount: 0,
		}
	}

	output.TokenCount = 10 // Placeholder for token count

	return output
}
