package handlers_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.iain.rocks/alectryon/api/handlers"
	"go.iain.rocks/alectryon/api/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestConvertModelToResponse(t *testing.T) {

	inputModel := models.InputModel{
		ID:        bson.NewObjectID(),
		Type:      models.InputTypeTelegramBot,
		Active:    true,
		CreatedAt: time.Now(),
	}

	outputResponse := handlers.ConvertModelToResponse(inputModel)

	assert.Equal(t, inputModel.ID.String(), outputResponse.Id)
	assert.Equal(t, string(inputModel.Type), outputResponse.Type)
	assert.Equal(t, inputModel.Active, outputResponse.Active)
	assert.Equal(t, inputModel.CreatedAt.Format(time.RFC3339), outputResponse.CreatedAt)
}
