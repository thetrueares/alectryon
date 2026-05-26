package vendor

import (
	"context"

	ollama "github.com/ollama/ollama/api"
	"go.iain.rocks/alectryon/api/engine"
	"go.iain.rocks/alectryon/api/entities"
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
	var output engine.Output
	var messages []ollama.Message

	for _, history := range input.History {
		messages = appendOllamaChatMessages(messages, history)
	}

	last := ollama.Message{
		Role:    "user",
		Content: input.Text,
	}
	messages = append(messages, last)

	request := &ollama.ChatRequest{
		Model:    "gemma4",
		Messages: messages,
		Stream:   &stream,
	}

	err := oa.client.Chat(context.TODO(), request, func(resp ollama.ChatResponse) error {
		output.Text += resp.Message.Content
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

func appendOllamaChatMessages(chatMessages []ollama.Message, history entities.HistoryEntity) []ollama.Message {

	var firstRole string
	var secondRole string
	if history.Direction == "inward" {
		firstRole = "user"
		secondRole = "assistant"
	} else {
		firstRole = "assistant"
		secondRole = "user"
	}

	first := ollama.Message{
		Role:    firstRole,
		Content: history.Message,
	}
	chatMessages = append(chatMessages, first)

	if history.Response != "" {
		second := ollama.Message{
			Role:    secondRole,
			Content: history.Response,
		}

		chatMessages = append(chatMessages, second)
	}

	return chatMessages
}
