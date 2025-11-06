package handler

import (
	"fmt"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type InfoCommandHandler struct {
	bot *tg.BotAPI
}

func NewInfoCommandHandler(bot *tg.BotAPI) *InfoCommandHandler {
	return &InfoCommandHandler{
		bot: bot,
	}
}

func (h *InfoCommandHandler) HandleCallback(callback *tg.CallbackQuery) {}

func (h *InfoCommandHandler) HandleCommand(callback *tg.Message) {
	h.bot.Send(tg.NewMessage(callback.Chat.ID, "Информация:\n"+
		"\n"+
		fmt.Sprintf("User ID: %d\n", callback.From.ID)+
		fmt.Sprintf("User name: %s\n", callback.From.UserName)+
		fmt.Sprintf("First name: %s\n", callback.From.FirstName)+
		fmt.Sprintf("Last name: %s\n", callback.From.LastName)+
		fmt.Sprintf("Language: %s\n", callback.From.LanguageCode),
	))
}
