package handler

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UnknownCommandHandler struct {
	bot *tg.BotAPI
}

func NewUnknownCommandHandler(bot *tg.BotAPI) *UnknownCommandHandler {
	return &UnknownCommandHandler{
		bot: bot,
	}
}

func (h *UnknownCommandHandler) HandleCallback(callback *tg.CallbackQuery) {}

func (h *UnknownCommandHandler) HandleCommand(callback *tg.Message) {
	h.bot.Send(tg.NewMessage(callback.Chat.ID, "❓ Неизвестная команда\n"+
		"\n"+
		"Переформулируйте свой вопрос без символа / или воспользуйтесь встроенной помощью:\n"+
		"/help - справочная информация\n"))
}
