package engine_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.iain.rocks/alectryon/backend/engine"
	"go.iain.rocks/alectryon/backend/entities"
)

func TestConvertTaskResponseToTask(t *testing.T) {
	taskResponse := engine.TaskResponse{
		Type:                entities.ProvideInformationTask,
		Description:         "Tell someone the time",
		ID:                  "6a12aeb6492aa2fc2671836a",
		RequiredInformation: map[string]string{},
	}

	taskEntity := engine.ConvertTaskResponseToTask(taskResponse)

	assert.NotEqual(t, 0, taskEntity.ID)
	assert.Equal(t, taskResponse.Type, taskEntity.Type)
	assert.Equal(t, taskResponse.Description, taskEntity.Description)
}
