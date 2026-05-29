package channels

import (
	"log"

	"go.iain.rocks/alectryon/backend/engine"
	"go.iain.rocks/alectryon/backend/entities"
)

func StartChannels(
	repository *entities.ChannelRepository,
	historyRepository *entities.HistoryRepository,
	userRepository *entities.UserRepository,
	engine engine.EngineInterface,
) {
	inputModels, err := repository.GetAll()

	if err != nil {
		panic(err.Error())
	}

	for _, input := range inputModels {
		if input.Type == entities.ChannelTypeTelegramBot {
			go func() {
				err := StartTelegramBot(input, historyRepository, userRepository, engine)
				if err != nil {
					log.Println(err.Error())
				}
			}()
		}
	}
}
