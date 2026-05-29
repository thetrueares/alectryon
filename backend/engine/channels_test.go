package engine_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.iain.rocks/alectryon/backend/engine"
	"go.iain.rocks/alectryon/backend/entities"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestCreateUserEntityFromInputMessage(t *testing.T) {
	inputMessage := engine.InputMessage{
		Channel: &entities.ChannelEntity{
			ID:   bson.NewObjectID(),
			Type: entities.ChannelTypeTelegramBot,
		},
		User: engine.InputUser{
			ChannelUserID: "third-id-one",
			Name:          "Iain Cambridge",
		},
	}
	userEntity := engine.CreateUserEntityFromInputMessage(inputMessage)

	assert.NotEqual(t, 0, userEntity.ID)
	assert.Equal(t, "Iain Cambridge", userEntity.Name)
	assert.Equal(t, "third-id-one", userEntity.UserChannels[0].UserID)
	assert.Equal(t, inputMessage.Channel.ID, userEntity.UserChannels[0].ChannelID)
	assert.Equal(t, entities.ChannelTypeTelegramBot, userEntity.UserChannels[0].ChannelType)
}
