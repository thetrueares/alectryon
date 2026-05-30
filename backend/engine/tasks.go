package engine

import (
	"go.iain.rocks/alectryon/backend/entities"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ActionType string

const (
	NewTaskAction     ActionType = "new_task"
	ResumedTaskAction ActionType = "resumed_task"
	ChatMessageAction ActionType = "chat_message"
)

type TaskResponse struct {
	ID                  string                                              `json:"id"`
	RequiredInformation map[string]entities.EmbeddedRequiredInformationData `json:"required_information"`
	Description         string                                              `json:"description"`
	Type                entities.TaskType                                   `json:"type"`
}

func ConvertTaskResponseToTask(taskResponse TaskResponse) *entities.TaskEntity {
	return &entities.TaskEntity{
		ID:                  bson.NewObjectID(),
		Type:                taskResponse.Type,
		Description:         taskResponse.Description,
		RequiredInformation: taskResponse.RequiredInformation,
	}
}

func AppendTaskWorkOutput(task *entities.TaskEntity, taskwork TaskWorkOutput) *entities.TaskEntity {
	if task == nil {
		// Handle the error gracefully or return early
		return task
	}
	embedded := entities.EmbeddedTaskWorkOutput{
		WorkDone: taskwork.WorkDone,
		Complete: taskwork.Complete,
		NextStep: taskwork.NextStep,
	}
	if task.TaskWorkOutput == nil {
		task.TaskWorkOutput = []entities.EmbeddedTaskWorkOutput{}
	}
	task.TaskWorkOutput = append(task.TaskWorkOutput, embedded)
	return task
}

type TaskWorkOutput struct {
	Complete bool         `json:"complete"`
	WorkDone string       `json:"work_done"`
	NextStep string       `json:"next_step"`
	Task     TaskResponse `json:"task"`
}
