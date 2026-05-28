package vendor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"slices"
	"time"

	ollama "github.com/ollama/ollama/api"
	"go.iain.rocks/alectryon/api/engine"
	"go.iain.rocks/alectryon/api/entities"
)

const MODEL_NAME = "gemma4"

type ReasonRequest struct {
	History   []ChatMessage `json:"history"`
	Latest    string        `json:"latest"`
	Timestamp string        `json:"timestamp"`
}

type ChatMessage struct {
	Role      string `json:"role"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

type ReasonResponse struct {
	History []ollama.Message  `json:"history"`
	Type    engine.ActionType `json:"type"`
}

func NewOllama() *Ollama {
	client, _ := ollama.ClientFromEnvironment()

	return &Ollama{client: client}
}

type Ollama struct {
	client *ollama.Client
}

func (oa Ollama) Process(input engine.ReasonResponse) engine.Output {
	stream := false
	var output engine.Output

	var messages []ollama.Message

	for _, msg := range input.History {
		msgObj := ollama.Message{
			Role:    msg.Role,
			Content: msg.Content,
		}
		messages = append(messages, msgObj)
	}

	last := ollama.Message{
		Role:    "user",
		Content: input.Latest,
	}
	messages = append(messages, last)
	request := &ollama.ChatRequest{
		Model:    MODEL_NAME,
		Messages: messages,
		Stream:   &stream,
		Think:    &ollama.ThinkValue{Value: "high"},
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

func (oa Ollama) Reason(input engine.Input) *engine.ReasonResponse {
	basePrompt := `With the following json body {history: [{role: \"user\", content: \"message\", timestamp: "2006-01-02T15:04:05Z07:00"}], latest: \"current_message\", timestamp: "2006-01-02T15:04:05Z07:00"}. 
    The response must say if it's a new chat, a resumed chat, a task, or a generate request. And if it's a resumed chat the history that is related to the chat must be returned in the history. Otherwise, the history is to be empty.
	For the it to be a new resumed chat the latest message must be related to the the history chat in subject.
	A task is something that is meant to be done
	The response must just be a json response with the body and no markdown {type: "resumed_chat|new_chat|task|generate", history: [{role: \"user\", content: \"message\"}], latest: "latest_message"}.
	The request payload is %s`

	messages := generateHistory(input.History)
	reasonRequest := ReasonRequest{
		History:   messages,
		Latest:    input.Text,
		Timestamp: time.Now().Format(time.RFC3339),
	}
	encodedStruct, err := json.Marshal(reasonRequest)

	if err != nil {
		log.Printf("[Ollama] json encoding failure \"%s\"\r\n", err.Error())

		return nil
	}

	prompt := fmt.Sprintf(basePrompt, string(encodedStruct))

	log.Printf("[Ollama] prompt \"%s\"\r\n", prompt)

	stream := false
	var output string
	generateRequest := &ollama.GenerateRequest{
		Model:  MODEL_NAME,
		Prompt: prompt,
		Stream: &stream,
	}
	clientErr := oa.client.Generate(context.TODO(), generateRequest, func(resp ollama.GenerateResponse) error {
		output += resp.Response

		return nil
	})

	if clientErr != nil {
		log.Printf("[Ollama] error: \"%s\"\r\n", clientErr.Error())

		return nil
	}

	log.Printf("[Ollama] reason response: %s\r\n", output)

	var reasonResp engine.ReasonResponse

	err = json.Unmarshal([]byte(output), &reasonResp)

	if err != nil {
		log.Printf("[Ollama] error: \"%s\"\r\n", err.Error())
	}

	return &reasonResp
}

func generateHistory(historyEntities []entities.HistoryEntity) []ChatMessage {

	var messages []ChatMessage
	cp := historyEntities
	slices.Reverse(cp)
	for _, history := range cp {
		messages = appendOllamaChatMessages(messages, history)
	}

	return messages
}

func appendOllamaChatMessages(chatMessages []ChatMessage, history entities.HistoryEntity) []ChatMessage {

	var firstRole string
	var secondRole string
	if history.Direction == "inward" {
		firstRole = "user"
		secondRole = "assistant"
	} else {
		firstRole = "assistant"
		secondRole = "user"
	}

	first := ChatMessage{
		Role:      firstRole,
		Content:   history.Message,
		Timestamp: history.CreatedAt.Format(time.RFC3339),
	}
	chatMessages = append(chatMessages, first)

	if history.Response != "" {
		second := ChatMessage{
			Role:      secondRole,
			Content:   history.Response,
			Timestamp: history.CreatedAt.Format(time.RFC3339),
		}

		chatMessages = append(chatMessages, second)
	}

	return chatMessages
}
