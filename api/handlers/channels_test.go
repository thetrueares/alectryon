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

	inputModel := models.ChannelEntity{
		ID:        bson.NewObjectID(),
		Name:      "Botty",
		Type:      models.InputTypeTelegramBot,
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	outputResponse := handlers.ConvertModelToResponse(inputModel)

	assert.Equal(t, inputModel.ID.Hex(), outputResponse.Id)
	assert.Equal(t, string(inputModel.Type), outputResponse.Type)
	assert.Equal(t, inputModel.Active, outputResponse.Active)
	assert.Equal(t, inputModel.CreatedAt.Format(time.RFC3339), outputResponse.CreatedAt)
	assert.Equal(t, inputModel.UpdatedAt.Format(time.RFC3339), outputResponse.UpdatedAt)
	assert.Equal(t, inputModel.Name, outputResponse.Name)
}

func TestConvertCreateRequestToModel(t *testing.T) {

	createRequest := handlers.ChannelCreateRequest{
		Name:   "Name",
		Type:   string(models.ChannelTypeTelegramBot),
		Active: true,
	}

	output := handlers.ConvertCreateRequestToModel(createRequest)

	assert.Equal(t, createRequest.Name, output.Name)
	assert.Equal(t, createRequest.Type, string(output.Type))
	assert.Equal(t, createRequest.Active, output.Active)
}
