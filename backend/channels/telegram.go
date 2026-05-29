package channels

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/go-telegram/bot"
	telegrammodels "github.com/go-telegram/bot/models"
	"go.iain.rocks/alectryon/backend/engine"
	"go.iain.rocks/alectryon/backend/entities"
)

// StartTelegramBot initializes and starts a Telegram bot based on the provided input model.
func StartTelegramBot(
	channel *entities.ChannelEntity,
	inputChan chan engine.InputMessage,
) error {
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

	senderHandler := make(chan engine.OutputMessage)
	handler := NewTelegramHandler(channel, inputChan, senderHandler)

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
	go sendMessage(context.Background(), b, senderHandler)

	return nil
}

func NewTelegramHandler(
	channelEntity *entities.ChannelEntity,
	inputChan chan engine.InputMessage,
	senderHandler chan engine.OutputMessage,
) *TelegramHandler {
	return &TelegramHandler{
		channelEntity: channelEntity,
		inputChan:     inputChan,
		senderHandler: senderHandler,
	}
}

type TelegramHandler struct {
	repository     *entities.HistoryRepository
	userRepository *entities.UserRepository
	channelEntity  *entities.ChannelEntity
	inputChan      chan engine.InputMessage
	senderHandler  chan engine.OutputMessage
}

func (th TelegramHandler) handle(ctx context.Context, b *bot.Bot, update *telegrammodels.Update) {
	if update.Message == nil {
		return
	}

	if update.Message.Text == "" {
		return
	}

	sender := "unknown"
	if update.Message.From != nil {
		sender = update.Message.From.Username
		if sender == "" {
			sender = fmt.Sprintf("%s %s", update.Message.From.FirstName, update.Message.From.LastName)
		}
	}

	log.Printf("[Telegram] Received message from %s: %s", sender, update.Message.Text)
	channelUserId := strconv.FormatInt(update.Message.From.ID, 10)
	inputMessage := engine.InputMessage{
		Channel:       th.channelEntity,
		Content:       update.Message.Text,
		SenderHandler: th.senderHandler,
		User: engine.InputUser{
			ChannelUserID: channelUserId,
			ChannelChatID: strconv.FormatInt(update.Message.Chat.ID, 10),
			Name:          sender,
		},
	}
	th.inputChan <- inputMessage
}

func sendMessage(ctx context.Context, b *bot.Bot, outputMessage chan engine.OutputMessage) {
	for message := range outputMessage {
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: message.InputUser.ChannelChatID,
			Text:   message.Content,
		})
		log.Printf("[Telegram] Sent message to %s: %s", message.InputUser.Name, message.Content)
		if err != nil {
			log.Printf("[Telegram] Failed to send message to %s: %v", message.InputUser.ChannelChatID, err)
		}
	}
}
