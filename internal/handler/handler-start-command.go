package handler

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type StartCommandHandler struct {
	BaseDependencies
}

func NewStartCommandHandler(baseDeps BaseDependencies) *StartCommandHandler {
	return &StartCommandHandler{BaseDependencies: baseDeps}
}

func (h *StartCommandHandler) HandleCommand(message *tg.Message) {
	h.SendTemplate(message.Chat.ID, "cmd-start.html", nil)
}
