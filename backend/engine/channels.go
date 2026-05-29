package engine

import (
	"go.iain.rocks/alectryon/backend/entities"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type InputMessage struct {
	Channel       *entities.ChannelEntity
	Content       string
	SenderHandler chan OutputMessage
	User          InputUser
}

type InputUser struct {
	ChannelUserID string
	ChannelChatID string
	Name          string
}

type OutputMessage struct {
	Content   string
	InputUser InputUser
}

func InputHandler(
	inputChan chan InputMessage,
	historyRepository *entities.HistoryRepository,
	userRepository *entities.UserRepository,
	ai EngineInterface,
) {
	for message := range inputChan {
		userEntity, err := userRepository.FindByChannelSender(message.Channel.Type, message.User.ChannelUserID)

		if err != nil {
			userEntity = createUserEntityFromInputMessage(message)
			userRepository.Save(userEntity)
		}

		history := entities.NewInwardMessage(userEntity, message.Content)

		aiOutput := ai.Process(Input{Text: message.Content, User: userEntity})
		history.Response = aiOutput.Text
		history.Task = entities.EmbeddedTask{
			ID:          bson.ObjectID([]byte(aiOutput.Task.ID)),
			Description: aiOutput.Task.Description,
		}
		historyRepository.Save(history)

		outputMessage := OutputMessage{
			Content:   aiOutput.Text,
			InputUser: message.User,
		}
		message.SenderHandler <- outputMessage
	}
}

func createUserEntityFromInputMessage(inputMessage InputMessage) *entities.UserEntity {

	return &entities.UserEntity{
		ID:   bson.NewObjectID(),
		Name: inputMessage.User.Name,
		UserChannels: []entities.UserChannel{
			{
				ChannelID:   inputMessage.Channel.ID,
				ChannelType: inputMessage.Channel.Type,
				UserID:      inputMessage.User.ChannelUserID,
			},
		},
	}
}
