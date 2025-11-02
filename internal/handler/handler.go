package handler

import (
	"fmt"
	"runtime/debug"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type (
	UpdateHandler interface {
		HandleCallback(callback *tg.CallbackQuery)
		HandleCommand(callback *tg.Message)
	}

	Handler struct {
		bot *tg.BotAPI
		log *logrus.Logger

		defaultUpdateHandler UpdateHandler
		// ...
	}
)

func NewHandler(bot *tg.BotAPI, log *logrus.Logger) *Handler {
	return &Handler{
		bot: bot,
		log: log,

		defaultUpdateHandler: nil,
	}
}

func (c *Handler) Handle(update tg.Update) {
	defer func() {
		if panicValue := recover(); panicValue != nil {
			c.log.Fatalf("recovered from panic: %v\n%v", panicValue, string(debug.Stack()))
		}
	}()

	// Save here ...

	switch {
	case update.CallbackQuery != nil:
		c.handleCallback(update.CallbackQuery)
	case update.Message != nil:
		c.handleMessage(update.Message)
	}
}

func (c *Handler) handleCallback(callback *tg.CallbackQuery) {
	switch callback.Data {
	case "test":
		break
	default:
		c.log.Warnf("Handler.handleCallback: unknown domain - %s", callback.Data)
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
		c.log.Warnf("Handler.handleCommand: unknown command - %s", message.Command())
	}
}

func (c *Handler) showWarningUnknownCommand(inputMessage *tg.Message) {
	c.log.Errorf("handler - unknown command - %v", inputMessage.Text)

	outputMessage := tg.NewMessage(inputMessage.Chat.ID, fmt.Sprintf("Unknown command - %s", inputMessage.Command()))

	_, err := c.bot.Send(outputMessage)
	if err != nil {
		c.log.Errorf("Handler.showWarningUnknownCommand: error sending reply message to chat - %v", err)
	}
}
