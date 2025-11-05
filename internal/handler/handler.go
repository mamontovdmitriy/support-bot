package handler

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
	"support-bot/config"
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
		cfg      config.TG
		bot      *tg.BotAPI
		services *service.Services

		defaultUpdateHandler UpdateHandler
		// unknownCommandHandler UpdateHandler
		// ...
	}
)

func NewHandler(cfg config.TG, bot *tg.BotAPI, services *service.Services) *Handler {
	return &Handler{
		cfg:      cfg,
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
	// case update.ChannelPost != nil:
	// c.handleChannelPost(update.ChannelPost)
	case update.CallbackQuery != nil:
		c.handleCallback(update.CallbackQuery)
	case update.Message != nil:
		c.handleMessage(update.Message)
	case update.EditedMessage != nil:
		c.handleMessage(update.EditedMessage)
	}
}

func (c *Handler) saveMessageUpdate(update tg.Update) {
	content, err := json.Marshal(update)
	if err != nil {
		c.services.Log.Warnf("can not marshal to json: %v", update)
		return
	}
	c.services.MessageUpdate.Create(entity.MessageUpdate{
		Id:      update.UpdateID,
		Message: string(content),
	})
}

func (c *Handler) handleCallback(callback *tg.CallbackQuery) {
	switch callback.Data {
	case "test":

	default:
		c.services.Log.Warnf("Handler.handleCallback: unknown domain - %s", callback.Data)
	}
}

func (c *Handler) handleMessage(message *tg.Message) {
	if !message.IsCommand() {
		c.proxyMessage(message)
		return
	}

	switch message.Command() {
	case "start":
	case "help":
	case "test":
		outputMessage := tg.NewMessage(message.Chat.ID, fmt.Sprintf("Handle command: %s", message.Command()))
		c.bot.Send(outputMessage)
	default:
		c.services.Log.Warnf("Handler.handleCommand: unknown command - %s", message.Command())
	}
}
