package channels_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.iain.rocks/alectryon/api/channels"
	"go.iain.rocks/alectryon/api/engine"
	"go.iain.rocks/alectryon/api/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type dummyAi struct{}

func (d dummyAi) Process(in engine.Input) engine.Output {
	return engine.Output{Text: in.Text}
}

func TestStartTelegramBot_InvalidType(t *testing.T) {
	input := models.InputModel{
		Type: models.InputTypeSlackBot,
	}
	err := channels.StartTelegramBot(input, nil, dummyAi{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid input type")
}

func TestStartTelegramBot_MissingToken(t *testing.T) {
	input := models.InputModel{
		Type:    models.InputTypeTelegramBot,
		Options: map[string]any{},
	}
	err := channels.StartTelegramBot(input, nil, dummyAi{})
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
	err := channels.StartTelegramBot(input, nil, dummyAi{})
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
	err := channels.StartTelegramBot(input, nil, dummyAi{})
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
	err := channels.StartTelegramBot(input, nil, dummyAi{})
	assert.Error(t, err)
}
