package handler

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type HelpCommandHandler struct {
	bot *tg.BotAPI
}

func NewHelpCommandHandler(bot *tg.BotAPI) *HelpCommandHandler {
	return &HelpCommandHandler{
		bot: bot,
	}
}

func (h *HelpCommandHandler) HandleCallback(callback *tg.CallbackQuery) {}

func (h *HelpCommandHandler) HandleCommand(callback *tg.Message) {
	h.bot.Send(tg.NewMessage(callback.Chat.ID, "Список доступных команд:\n"+
		"\n"+
		"/start - рекомендации для начала работы\n"+
		"/help - справочная информация\n"))
}
