package handler

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type HelpCommandHandler struct {
	BaseDependencies
}

func NewHelpCommandHandler(baseDeps BaseDependencies) *HelpCommandHandler {
	return &HelpCommandHandler{BaseDependencies: baseDeps}
}

func (h *HelpCommandHandler) HandleCommand(message *tg.Message) {
	h.SendTemplate(message.Chat.ID, "cmd-help.html", nil)
}
