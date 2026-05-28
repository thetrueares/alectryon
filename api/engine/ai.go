package engine

import (
	"go.iain.rocks/alectryon/api/entities"
)

type Input struct {
	Text    string
	History []entities.HistoryEntity
	User    *entities.UserEntity
}

type Output struct {
	Text       string
	TokenCount int
}

type ReasonResponse struct {
	Action  ActionType
	History []ChatMessage
	Latest  string
}

type ChatMessage struct {
	Role    string
	Content string
}

type ActionType string

const (
	NewChatAction     ActionType = "new_chat"
	ResumedChatAction ActionType = "resumed_chat"
	GenerateAction    ActionType = "generate"
	TaskAction        ActionType = "task"
)

type EngineInterface interface {
	Process(in Input) Output
}

type Engine struct {
	ai                AiInterface
	historyRepository entities.HistoryRepository
}

func (e Engine) Process(in Input) Output {
	in.History, _ = e.historyRepository.GetLastTenForUser(in.User)

	resp := e.ai.Reason(in)

	return e.ai.Process(*resp)
}

func NewEngine(ai AiInterface, historyRepository entities.HistoryRepository) EngineInterface {
	return &Engine{ai: ai, historyRepository: historyRepository}
}

type AiInterface interface {
	Process(input ReasonResponse) Output
	Reason(input Input) *ReasonResponse
}

type SimpleAi struct{}

func (s SimpleAi) Process(in Input) Output {
	return Output{Text: "AI Response: " + in.Text}
}
