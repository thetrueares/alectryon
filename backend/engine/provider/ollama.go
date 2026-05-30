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
	output.Task = input.Task
	return output
}

func (oa Ollama) Reason(input engine.Input) *engine.ReasonResponse {
	basePrompt := `You are a personal assistant who has to reason what the latest message is and what the outcome of the message is. 
If it's a new task you are provide a description of the task.  
The related history must be directly related to the task. 
In some cases more information is to be requested to be able to fulfill the task.
With the following json body {history: [{role: "user", content: "message", timestamp: "2006-01-02T15:04:05Z07:00", task_id: "id"}], latest: "current_message", timestamp: "2006-01-02T15:04:05Z07:00"}.
The response must say if it's a new task, an existing task, or a generate request. And if it's a resumed task the history that is related to the task must be returned in the history. Otherwise, the history is to be empty.
For the it to be a resumed task the latest message must be related to the the history chat in subject.
If the task is new then don't return a task id otherwise use the task id for the related history messages. All the related history messages must have the same task id.
A task is something that is meant to be done. Some tasks are that something needs to be done, these are action verbs or needs. Some are provide information and these are when they ask questions.
The response must just be a json response with the body and no markdown {type: "resumed_task|new_task", history: [{role: "user", content: "message", task_id: "id"], latest: "latest_message", "expected_outcome": "expected_outcome", task:{"id": "id"," "type": "PERFORM_ACTION|PROVIDE_INFORMATION"," "description": "new_description", required_information: {}}}.
The request payload is %s`

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

	prompt := fmt.Sprintf(basePrompt, string(encodedStruct))

	stream := false
	var output string
	generateRequest := &ollama.GenerateRequest{
		Model:  MODEL_NAME,
		Prompt: prompt,
		Stream: &stream,
	}

	jsonRequest, err := json.Marshal(generateRequest)

	if err != nil {
		oa.logger.Error(fmt.Sprintf("[Ollama] json marshalling error: \"%s\"\r\n", err.Error()))

		return nil
	}

	oa.logger.Info(fmt.Sprintf("[Ollama] json request \"%s\"\r\n", string(jsonRequest)))

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
