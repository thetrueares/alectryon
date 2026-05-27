package channels

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/go-telegram/bot"
	telegrammodels "github.com/go-telegram/bot/models"
	"go.iain.rocks/alectryon/api/engine"
	"go.iain.rocks/alectryon/api/entities"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// StartTelegramBot initializes and starts a Telegram bot based on the provided input model.
func StartTelegramBot(channel entities.ChannelEntity, repository *entities.HistoryRepository, userRepository *entities.UserRepository, ai engine.EngineInterface) error {
	if channel.Type != entities.ChannelTypeTelegramBot {
		return fmt.Errorf("invalid input type for telegram: %s", channel.Type)
	}

	token, ok := channel.Options["bot_token"]
	if !ok {
		return errors.New("bot_token not found in input options")
	}

	tokenStr, ok := token.(string)
	if !ok {
		return errors.New("bot_token must be a string")
	}

	if tokenStr == "" {
		return errors.New("bot_token is empty")
	}

	handler := NewTelegramHandler(channel, repository, userRepository, ai)

	opts := []bot.Option{
		bot.WithDefaultHandler(handler.handle),
	}

	b, err := bot.New(tokenStr, opts...)
	if err != nil {
		return fmt.Errorf("failed to create telegram bot: %w", err)
	}

	log.Printf("Starting Telegram bot for input: %s (%s)", channel.Name, channel.ID.Hex())

	// Start the bot in a goroutine
	go b.Start(context.TODO())

	return nil
}

func NewTelegramHandler(channelEntity entities.ChannelEntity, repository *entities.HistoryRepository, userRepository *entities.UserRepository, ai engine.EngineInterface) *TelegramHandler {
	return &TelegramHandler{repository: repository, userRepository: userRepository, ai: ai, channelEntity: channelEntity}
}

type TelegramHandler struct {
	repository     *entities.HistoryRepository
	userRepository *entities.UserRepository
	ai             engine.EngineInterface
	channelEntity  entities.ChannelEntity
}

func (th TelegramHandler) handle(ctx context.Context, b *bot.Bot, update *telegrammodels.Update) {
	if update.Message == nil {
		return
	}

	sender := "unknown"
	if update.Message.From != nil {
		sender = update.Message.From.Username
		if sender == "" {
			sender = update.Message.From.FirstName
		}
	}

	if update.Message.Text == "" {
		return
	}

	userEntity, err := th.userRepository.FindByChannelSender(entities.ChannelTypeTelegramBot, strconv.FormatInt(update.Message.From.ID, 10))

	if err != nil {
		userEntity = createUserEntityFromTelegramSender(*update.Message.From, th.channelEntity)
		th.userRepository.Save(*userEntity)
	}

	history := entities.NewInwardMessage(userEntity, update.Message.Text)
	log.Printf("[Telegram] Received message from %s: %s", sender, update.Message.Text)

	aiOutput := th.ai.Process(engine.Input{Text: update.Message.Text})
	history.Response = aiOutput.Text
	th.repository.Save(history)

	th.sendMessage(ctx, b, update.Message.Chat.ID, sender, aiOutput.Text)
}

func (th TelegramHandler) sendMessage(ctx context.Context, b *bot.Bot, chatId int64, user, message string) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatId,
		Text:   message,
	})
	if err != nil {
		log.Printf("[Telegram] Failed to send message to %d: %v", chatId, err)
	}
}

func createUserEntityFromTelegramSender(telegramSender telegrammodels.User, channel entities.ChannelEntity) *entities.UserEntity {

	return &entities.UserEntity{
		ID:   bson.NewObjectID(),
		Name: telegramSender.FirstName,
		UserChannels: []entities.UserChannel{
			{
				ChannelID:   channel.ID,
				ChannelType: entities.ChannelTypeTelegramBot,
				UserID:      strconv.FormatInt(telegramSender.ID, 10),
			},
		},
	}
}
