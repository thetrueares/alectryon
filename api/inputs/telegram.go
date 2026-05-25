package inputs

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/go-telegram/bot"
	telegrammodels "github.com/go-telegram/bot/models"
	"go.iain.rocks/alectryon/api/models"
)

// StartTelegramBot initializes and starts a Telegram bot based on the provided input model.
func StartTelegramBot(input models.InputModel, repository *models.HistoryRepository) error {
	if input.Type != models.InputTypeTelegramBot {
		return fmt.Errorf("invalid input type for telegram: %s", input.Type)
	}

	token, ok := input.Options["bot_token"]
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

	handler := NewTelegramHandler(repository)

	opts := []bot.Option{
		bot.WithDefaultHandler(handler.handle),
	}

	b, err := bot.New(tokenStr, opts...)
	if err != nil {
		return fmt.Errorf("failed to create telegram bot: %w", err)
	}

	log.Printf("Starting Telegram bot for input: %s (%s)", input.Name, input.ID.Hex())

	// Start the bot in a goroutine
	go b.Start(context.TODO())

	return nil
}

func NewTelegramHandler(repository *models.HistoryRepository) *TelegramHandler {
	return &TelegramHandler{repository: repository}
}

type TelegramHandler struct {
	repository *models.HistoryRepository
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

	history := models.NewInwardMessage(sender, update.Message.Text)
	th.repository.Save(history)

	log.Printf("[Telegram] Received message from %s: %s", sender, update.Message.Text)

	history.Response = update.Message.Text
	th.repository.Save(history)
	th.sendMessage(ctx, b, update.Message.Chat.ID, sender, update.Message.Text)
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
