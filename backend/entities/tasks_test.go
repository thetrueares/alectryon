package entities_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.iain.rocks/alectryon/backend/entities"
)

func TestTaskRepository_FindById_InvalidID(t *testing.T) {
	tr := entities.NewTaskRepository(nil)
	task, err := tr.FindById("invalid-id")
	
	assert.Error(t, err)
	assert.Nil(t, task)
}
