package handler

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type StartCommandHandler struct {
	bot *tg.BotAPI
}

func NewStartCommandHandler(bot *tg.BotAPI) *StartCommandHandler {
	return &StartCommandHandler{
		bot: bot,
	}
}

func (h *StartCommandHandler) HandleCallback(callback *tg.CallbackQuery) {}

func (h *StartCommandHandler) HandleCommand(callback *tg.Message) {
	h.bot.Send(tg.NewMessage(callback.Chat.ID, "📬 Служба поддержки клиентов\n"+
		"\n"+
		"Здесь Вы можете проконсультироваться и получить ответы непосредственно от официальных представителей компании.\n"+
		"\n"+
		"Это легко и удобно! Просто напишите Ваш вопрос🙂\n"))
}
