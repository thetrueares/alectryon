package engine

import (
	"go.iain.rocks/alectryon/backend/entities"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.uber.org/zap"
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
	logger *zap.Logger,
) {
	for message := range inputChan {
		logger.Debug("Processing message", zap.Any("message", message.Content))
		userEntity, err := userRepository.FindByChannelSender(message.Channel.Type, message.User.ChannelUserID)

		if err != nil {
			logger.Debug("Creating new user from input message")
			userEntity = CreateUserEntityFromInputMessage(message)
			err = userRepository.Save(userEntity)
			if err != nil {
				logger.Error("Error saving user to database", zap.Field{
					Key:    "error",
					String: err.Error(),
				})
			}
		}

		history := entities.NewInwardMessage(userEntity, message.Content)

		aiOutput := ai.Process(Input{Text: message.Content, User: userEntity})
		history.Response = aiOutput.Text
		history.Task = entities.EmbeddedTask{
			ID:          bson.ObjectID([]byte(aiOutput.Task.ID)),
			Description: aiOutput.Task.Description,
		}
		err = historyRepository.Save(history)

		if err != nil {
			logger.Error("Error saving history to database", zap.Field{
				Key:    "error",
				String: err.Error(),
			})
		}
		outputMessage := OutputMessage{
			Content:   aiOutput.Text,
			InputUser: message.User,
		}
		logger.Debug("Sending output message", zap.Any("content", outputMessage.Content))
		message.SenderHandler <- outputMessage
	}
}

func CreateUserEntityFromInputMessage(inputMessage InputMessage) *entities.UserEntity {

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
