package channels

import (
	"log"

	"go.iain.rocks/alectryon/backend/engine"
	"go.iain.rocks/alectryon/backend/entities"
)

func StartChannels(
	repository *entities.ChannelRepository,
	messageHandler chan engine.InputMessage,
) {
	channelsModels, err := repository.GetAll()

	if err != nil {
		panic(err.Error())
	}

	for _, channel := range channelsModels {
		if channel.Type == entities.ChannelTypeTelegramBot {
			go func() {
				err := StartTelegramBot(&channel, messageHandler)
				if err != nil {
					log.Println(err.Error())
				}
			}()
		}
	}
}
