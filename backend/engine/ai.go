package engine

import (
	"log"

	"github.com/bytedance/gopkg/util/logger"
	"go.iain.rocks/alectryon/backend/entities"
)

type Input struct {
	Text    string
	History []entities.HistoryEntity
	User    *entities.UserEntity
}

type Output struct {
	Text       string
	TokenCount int
	Task       TaskResponse
}

type ReasonResponse struct {
	Type    ActionType    `json:"type"`
	History []ChatMessage `json:"history"`
	Latest  string        `json:"latest,omitempty"`
	Task    TaskResponse  `json:"task"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	TaskID  string `json:"task_id"`
}

type EngineInterface interface {
	Process(in Input) Output
}

type Engine struct {
	ai                AiInterface
	historyRepository *entities.HistoryRepository
	taskRepository    *entities.TaskRepository
}

func (e Engine) Process(in Input) Output {
	in.History, _ = e.historyRepository.GetLastTenForUser(in.User)

	resp := e.ai.Reason(in)

	if resp.Type == NewTaskAction {
		taskEntity := ConvertTaskResponseToTask(resp.Task)
		err := e.taskRepository.Save(taskEntity)
		if err != nil {
			logger.Error(err.Error())
		}
		resp.Task.ID = taskEntity.ID.Hex()
	}

	log.Print(resp.Type)

	return e.ai.Process(*resp)
}

func NewEngine(
	ai AiInterface,
	historyRepository *entities.HistoryRepository,
	taskRepository *entities.TaskRepository,
) EngineInterface {
	return &Engine{ai: ai, historyRepository: historyRepository, taskRepository: taskRepository}
}

type AiInterface interface {
	Process(input ReasonResponse) Output
	Reason(input Input) *ReasonResponse
}

type SimpleAi struct{}

func (s SimpleAi) Process(in Input) Output {
	return Output{Text: "AI Response: " + in.Text}
}
func (s SimpleAi) Reason(input Input) *ReasonResponse { return &ReasonResponse{} }
