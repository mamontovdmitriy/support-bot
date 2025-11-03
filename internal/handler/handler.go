package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime/debug"
	"support-bot/internal/entity"
	"support-bot/internal/service"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type (
	UpdateHandler interface {
		HandleCallback(callback *tg.CallbackQuery)
		HandleCommand(callback *tg.Message)
	}

	Handler struct {
		bot      *tg.BotAPI
		services *service.Services

		defaultUpdateHandler UpdateHandler
		// ...
	}
)

func NewHandler(bot *tg.BotAPI, services *service.Services) *Handler {
	return &Handler{
		bot:      bot,
		services: services,

		defaultUpdateHandler: nil,
	}
}

func (c *Handler) Handle(update tg.Update) {
	defer func() {
		if panicValue := recover(); panicValue != nil {
			c.services.Log.Fatalf("recovered from panic: %v\n%v", panicValue, string(debug.Stack()))
		}
	}()

	// Save all request here ...
	c.saveMessageUpdate(update)

	// Processing message or callback
	switch {
	case update.CallbackQuery != nil:
		c.handleCallback(update.CallbackQuery)
	case update.Message != nil:
		c.handleMessage(update.Message)
	}
}

func (c *Handler) saveMessageUpdate(update tg.Update) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	content, err := json.Marshal(update)
	if err != nil {
		c.services.Log.Warnf("can not marshal to json: %v", update)
		return
	}
	c.services.MessageUpdate.Create(ctx, entity.MessageUpdate{
		Id:      update.UpdateID,
		Message: string(content),
	})
}

func (c *Handler) handleCallback(callback *tg.CallbackQuery) {
	switch callback.Data {
	case "test":
		break
	default:
		c.services.Log.Warnf("Handler.handleCallback: unknown domain - %s", callback.Data)
	}
}

func (c *Handler) handleMessage(message *tg.Message) {
	if !message.IsCommand() {
		c.showWarningUnknownCommand(message)
		return
	}

	switch message.Command() {
	case "test":
		outputMessage := tg.NewMessage(message.Chat.ID, fmt.Sprintf("Handle command: %s", message.Command()))
		c.bot.Send(outputMessage)
		break
	default:
		c.services.Log.Warnf("Handler.handleCommand: unknown command - %s", message.Command())
	}
}

func (c *Handler) showWarningUnknownCommand(inputMessage *tg.Message) {
	c.services.Log.Errorf("handler - unknown command - %v", inputMessage.Text)

	outputMessage := tg.NewMessage(inputMessage.Chat.ID, fmt.Sprintf("Unknown command - %s", inputMessage.Command()))

	_, err := c.bot.Send(outputMessage)
	if err != nil {
		c.services.Log.Errorf("Handler.showWarningUnknownCommand: error sending reply message to chat - %v", err)
	}
}
