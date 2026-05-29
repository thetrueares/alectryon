package engine

import (
	"go.iain.rocks/alectryon/backend/entities"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ActionType string

const (
	NewTaskAction     ActionType = "new_task"
	ResumedTaskAction ActionType = "resumed_task"
)

type TaskResponse struct {
	ID                  string            `json:"id"`
	RequiredInformation map[string]string `json:"required_information"`
	Description         string            `json:"description"`
	Type                entities.TaskType `json:"type"`
}

func ConvertTaskResponseToTask(taskResponse TaskResponse) *entities.TaskEntity {
	return &entities.TaskEntity{
		ID:                  bson.NewObjectID(),
		Type:                taskResponse.Type,
		Description:         taskResponse.Description,
		RequiredInformation: taskResponse.RequiredInformation,
	}
}
