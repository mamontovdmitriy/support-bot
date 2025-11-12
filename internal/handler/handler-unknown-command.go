package handler

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UnknownCommandHandler struct {
	BaseDependencies
}

func NewUnknownCommandHandler(baseDeps BaseDependencies) *UnknownCommandHandler {
	return &UnknownCommandHandler{BaseDependencies: baseDeps}
}

func (h *UnknownCommandHandler) HandleCommand(message *tg.Message) {
	h.SendTemplate(message.Chat.ID, "cmd-unknown.html", nil)
}
