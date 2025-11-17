package handler

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"runtime/debug"

	"support-bot/config"
	"support-bot/internal/entity"
	"support-bot/internal/service"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type (
	CallbackHandler interface {
		HandleCallback(callback *tg.CallbackQuery)
	}
	ChatMemberHandler interface {
		HandleChatMember(chatMember *tg.ChatMemberUpdated)
	}
	CommandHandler interface {
		HandleCommand(message *tg.Message)
	}
	MessageHandler interface {
		CommandHandler
	}

	UpdateHandler interface {
		CallbackHandler
		ChatMemberHandler
		MessageHandler
	}

	BaseDependencies struct {
		cfg       config.TG
		bot       *tg.BotAPI
		services  *service.Services
		templates *template.Template
	}

	Handler struct {
		BaseDependencies

		defaultUpdateHandler UpdateHandler
		// Commands
		startCommandHandler   CommandHandler
		helpCommandHandler    CommandHandler
		infoCommandHandler    CommandHandler
		unknownCommandHandler CommandHandler
		// Chat member add/edit
		channelChatMemberHandler ChatMemberHandler
		groupChatMemberHandler   ChatMemberHandler
	}
)

//go:embed templates/*.html
var htmlFiles embed.FS

func NewHandler(cfg config.TG, bot *tg.BotAPI, services *service.Services) *Handler {
	templates, err := template.ParseFS(htmlFiles, "templates/*.html")
	if err != nil {
		services.Log.Fatal("Handler.NewHandler: error parsing templates - ", err)
	}

	baseDeps := BaseDependencies{
		cfg:       cfg,
		bot:       bot,
		services:  services,
		templates: templates,
	}

	return &Handler{
		BaseDependencies: baseDeps,

		defaultUpdateHandler: nil,
		// Commands
		startCommandHandler:   NewStartCommandHandler(baseDeps),
		helpCommandHandler:    NewHelpCommandHandler(baseDeps),
		infoCommandHandler:    NewInfoCommandHandler(baseDeps),
		unknownCommandHandler: NewUnknownCommandHandler(baseDeps),
		// Chat member add/edit
		channelChatMemberHandler: NewChannelMemberHandler(baseDeps),
		groupChatMemberHandler:   NewGroupChatMemberHandler(baseDeps),
	}
}

func (c *Handler) Handle(update tg.Update) {
	defer func() {
		if panicValue := recover(); panicValue != nil {
			c.services.Log.Fatalf("recovered from panic: %v\n%v", panicValue, string(debug.Stack()))
		}
	}()

	// Save all request in db
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
	case update.MyChatMember != nil:
		c.handleChatMember(update.MyChatMember)
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
		c.startCommandHandler.HandleCommand(message)
	case "help":
		c.helpCommandHandler.HandleCommand(message)
	case "info":
		c.infoCommandHandler.HandleCommand(message)
	default:
		c.services.Log.Warnf("Handler.handleCommand: unknown command - %s", message.Command())
		c.unknownCommandHandler.HandleCommand(message)
	}
}

func (c *Handler) handleChatMember(chatMember *tg.ChatMemberUpdated) {
	switch chatMember.Chat.Type {
	case "channel": // public for posts
		c.channelChatMemberHandler.HandleChatMember(chatMember)
	case "supergroup": // group for post`s comments
		c.groupChatMemberHandler.HandleChatMember(chatMember)
	}
}

func (bh *BaseDependencies) SendTemplate(chatID int64, tmplName string, data interface{}) {
	text, err := RenderTemplate(bh.templates, tmplName, data)
	if err != nil {
		bh.services.Log.Error(tmplName, ": error template - ", err)
	}

	_, err = bh.bot.Send(tg.MessageConfig{
		BaseChat: tg.BaseChat{
			ChatID:           chatID,
			ReplyToMessageID: 0,
		},
		Text:                  text,
		ParseMode:             "HTML",
		DisableWebPagePreview: false,
	})
	if err != nil {
		bh.services.Log.Error(tmplName, ": error sending - ", err)
		bh.bot.Send(tg.NewMessage(chatID, fmt.Sprintf("❗ERROR: %v", err)))
	}
}

func RenderTemplate(templates *template.Template, templateName string, data interface{}) (string, error) {
	var buf bytes.Buffer
	err := templates.ExecuteTemplate(&buf, templateName, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
