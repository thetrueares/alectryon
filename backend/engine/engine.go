package engine

import (
	"go.iain.rocks/alectryon/backend/entities"
	"go.uber.org/zap"
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

type ProcessedMessage struct {
	Task                 TaskResponse   `json:"task"`
	LatestTaskWorkOutput TaskWorkOutput `json:"latest_task_work_output"`
	LatestMessage        string         `json:"latest_message"`
}

type EngineInterface interface {
	Process(in Input) Output
}

type Engine struct {
	ai                AiInterface
	historyRepository *entities.HistoryRepository
	taskRepository    *entities.TaskRepository
	logger            *zap.Logger
}

func (e Engine) Process(in Input) Output {
	e.logger.Info("Processing input", zap.String("input", in.Text))

	in.History, _ = e.historyRepository.GetLastTenForUser(in.User)

	resp := e.ai.Reason(in)

	if resp.Type == NewTaskAction {
		taskEntity := ConvertTaskResponseToTask(resp.Task)
		err := e.taskRepository.Save(taskEntity)

		if err != nil {
			e.logger.Error("Error saving task to database", zap.Any("error", err.Error()))
		}

		resp.Task.ID = taskEntity.ID.Hex()
	}

	e.logger.Info("Processing output", zap.String("res", in.Text))

	outcome := ProcessedMessage{
		Task:                 resp.Task,
		LatestTaskWorkOutput: TaskWorkOutput{},
		LatestMessage:        in.Text,
	}
	return e.ai.Process(outcome)
}

func NewEngine(
	ai AiInterface,
	historyRepository *entities.HistoryRepository,
	taskRepository *entities.TaskRepository,
	logger *zap.Logger,
) EngineInterface {
	return &Engine{ai: ai, historyRepository: historyRepository, taskRepository: taskRepository, logger: logger}
}

type AiInterface interface {
	Process(input ProcessedMessage) Output
	Reason(input Input) *ReasonResponse
}
