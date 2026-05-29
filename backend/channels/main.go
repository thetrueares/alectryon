package channels

import (
	"fmt"
	"log"

	"go.iain.rocks/alectryon/backend/engine"
	"go.iain.rocks/alectryon/backend/entities"
	"go.uber.org/zap"
)

func StartChannels(
	repository *entities.ChannelRepository,
	messageHandler chan engine.InputMessage,
	logger *zap.Logger,
) {
	channelsModels, err := repository.GetAll()

	if err != nil {
		panic(err.Error())
	}

	for _, channel := range channelsModels {
		if channel.Type == entities.ChannelTypeTelegramBot {
			go func() {
				logger.Info(fmt.Sprintf("Starting telegram bot for channel %s", channel.ID))
				err := StartTelegramBot(&channel, messageHandler, logger)
				if err != nil {
					log.Println(err.Error())
				}
			}()
		}
	}
}
