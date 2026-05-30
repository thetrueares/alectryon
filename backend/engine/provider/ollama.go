package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"time"

	ollama "github.com/ollama/ollama/api"
	"go.iain.rocks/alectryon/backend/engine"
	"go.iain.rocks/alectryon/backend/entities"
	"go.uber.org/zap"
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
	TaskID    string `json:"task_id"`
}

type ReasonResponse struct {
	History []ollama.Message  `json:"history"`
	Type    engine.ActionType `json:"type"`
}

func NewOllama(logger *zap.Logger) *Ollama {
	client, _ := ollama.ClientFromEnvironment()

	return &Ollama{client: client, logger: logger}
}

type Ollama struct {
	client *ollama.Client
	logger *zap.Logger
}

func (oa Ollama) Process(input engine.ProcessedMessage) engine.Output {

	var output engine.Output
	encodedStruct, err := json.Marshal(input)

	if err != nil {
		oa.logger.Error(fmt.Sprintf("[Ollama] json encoding failure \"%s\"\r\n", err.Error()))
		output.Text = "An error happened during json encoding"
		return output
	}

	prompt := fmt.Sprintf(engine.PromptOutput, string(encodedStruct))
	generateRequest := newGenerateRequest(prompt)
	jsonRequest, err := json.Marshal(generateRequest)

	if err != nil {
		oa.logger.Error(fmt.Sprintf("[Ollama] json marshalling error: \"%s\"\r\n", err.Error()))
		output.Text = "An error happened during json encoding for the json request to ollama"
		return output
	}

	oa.logger.Info(fmt.Sprintf("[Ollama] json request \"%s\"\r\n", string(jsonRequest)))
	err = oa.client.Generate(context.TODO(), generateRequest, func(resp ollama.GenerateResponse) error {
		output.Text += resp.Response
		return nil
	})

	oa.logger.Info(fmt.Sprintf("[Ollama] process response: %s\r\n", output.Text))
	if err != nil {
		return engine.Output{
			Text:       "Ollama error: " + err.Error(),
			TokenCount: 0,
		}
	}

	output.TokenCount = 10 // Placeholder for token count
	output.Task = input.Task
	return output
}

func (oa Ollama) AnalyseTask(input *engine.ReasonResponse) engine.TaskWorkOutput {
	encodedStruct, err := json.Marshal(input)

	if err != nil {
		oa.logger.Error(fmt.Sprintf("[Ollama] json encoding failure \"%s\"\r\n", err.Error()))
		return engine.TaskWorkOutput{}
	}

	prompt := fmt.Sprintf(engine.PromptTaskNextStep, string(encodedStruct))
	generateRequest := newGenerateRequest(prompt)
	jsonRequest, err := json.Marshal(generateRequest)

	if err != nil {
		oa.logger.Error(fmt.Sprintf("[Ollama] json marshalling error: \"%s\"\r\n", err.Error()))
		return engine.TaskWorkOutput{}
	}

	var output string
	oa.logger.Info(fmt.Sprintf("[Ollama] json request \"%s\"\r\n", string(jsonRequest)))
	err = oa.client.Generate(context.TODO(), generateRequest, func(resp ollama.GenerateResponse) error {
		output = resp.Response
		return nil
	})

	oa.logger.Info(fmt.Sprintf("[Ollama] analyse task response: %s\r\n", output))
	var reasonResp engine.TaskWorkOutput

	err = json.Unmarshal([]byte(output), &reasonResp)

	if err != nil {
		oa.logger.Error(fmt.Sprintf("[Ollama] error: \"%s\"\r\n", err.Error()))
	}

	return reasonResp
}

func (oa Ollama) Reason(input engine.Input) *engine.ReasonResponse {
	messages := generateHistory(input.History)
	reasonRequest := ReasonRequest{
		History:   messages,
		Latest:    input.Text,
		Timestamp: time.Now().Format(time.RFC3339),
	}
	encodedStruct, err := json.Marshal(reasonRequest)

	if err != nil {
		oa.logger.Error(fmt.Sprintf("[Ollama] json encoding failure \"%s\"\r\n", err.Error()))

		return nil
	}

	prompt := fmt.Sprintf(engine.PromptTaskIdentification, string(encodedStruct))
	generateRequest := newGenerateRequest(prompt)
	jsonRequest, err := json.Marshal(generateRequest)

	if err != nil {
		oa.logger.Error(fmt.Sprintf("[Ollama] json marshalling error: \"%s\"\r\n", err.Error()))

		return nil
	}

	oa.logger.Info(fmt.Sprintf("[Ollama] json request \"%s\"\r\n", string(jsonRequest)))

	var output string
	clientErr := oa.client.Generate(context.TODO(), generateRequest, func(resp ollama.GenerateResponse) error {
		output += resp.Response

		return nil
	})

	if clientErr != nil {
		oa.logger.Error(fmt.Sprintf("[Ollama] error: \"%s\"\r\n", clientErr.Error()))

		return nil
	}

	oa.logger.Info(fmt.Sprintf("[Ollama] reason response: %s\r\n", output))

	var reasonResp engine.ReasonResponse

	err = json.Unmarshal([]byte(output), &reasonResp)

	if err != nil {
		oa.logger.Error(fmt.Sprintf("[Ollama] error: \"%s\"\r\n", err.Error()))
	}

	return &reasonResp
}

func newGenerateRequest(prompt string) *ollama.GenerateRequest {

	return &ollama.GenerateRequest{
		Model:  MODEL_NAME,
		Prompt: prompt,
		Stream: new(false),
		Think:  &ollama.ThinkValue{"max"},
	}
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
		TaskID:    history.Task.ID.Hex(),
	}
	chatMessages = append(chatMessages, first)

	if history.Response != "" {
		second := ChatMessage{
			Role:      secondRole,
			Content:   history.Response,
			Timestamp: history.CreatedAt.Format(time.RFC3339),
			TaskID:    history.Task.ID.Hex(),
		}

		chatMessages = append(chatMessages, second)
	}

	return chatMessages
}
