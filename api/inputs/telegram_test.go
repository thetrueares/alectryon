package inputs_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.iain.rocks/alectryon/api/inputs"
	"go.iain.rocks/alectryon/api/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestStartTelegramBot_InvalidType(t *testing.T) {
	input := models.InputModel{
		Type: models.InputTypeSlackBot,
	}
	err := inputs.StartTelegramBot(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid input type")
}

func TestStartTelegramBot_MissingToken(t *testing.T) {
	input := models.InputModel{
		Type:    models.InputTypeTelegramBot,
		Options: map[string]any{},
	}
	err := inputs.StartTelegramBot(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bot_token not found")
}

func TestStartTelegramBot_InvalidTokenFormat(t *testing.T) {
	input := models.InputModel{
		Type: models.InputTypeTelegramBot,
		Options: map[string]any{
			"bot_token": 12345,
		},
	}
	err := inputs.StartTelegramBot(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bot_token must be a string")
}

func TestStartTelegramBot_EmptyToken(t *testing.T) {
	input := models.InputModel{
		Type: models.InputTypeTelegramBot,
		Options: map[string]any{
			"bot_token": "",
		},
	}
	err := inputs.StartTelegramBot(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bot_token is empty")
}

func TestStartTelegramBot_InvalidTokenValue(t *testing.T) {
	// This will fail because bot.New validates the token format (usually)
	input := models.InputModel{
		ID:   bson.NewObjectID(),
		Name: "Test Bot",
		Type: models.InputTypeTelegramBot,
		Options: map[string]any{
			"bot_token": "invalid-token",
		},
	}
	err := inputs.StartTelegramBot(input)
	assert.Error(t, err)
}
