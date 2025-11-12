package handler

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ChannelChatMemberHandler struct {
	BaseDependencies
}

func NewChannelMemberHandler(baseDeps BaseDependencies) *ChannelChatMemberHandler {
	return &ChannelChatMemberHandler{BaseDependencies: baseDeps}
}

func (h *ChannelChatMemberHandler) HandleChatMember(chatMemberUpdated *tg.ChatMemberUpdated) {
	newMember := chatMemberUpdated.NewChatMember

	if !newMember.User.IsBot {
		return
	}

	var criteria = []Criterium{
		Criterium{"Статус", func(chm tg.ChatMember) (string, bool) {
			return chm.Status, chm.IsAdministrator()
		}},
		Criterium{"Управление чатом", func(chm tg.ChatMember) (string, bool) {
			result := "Нет"
			if chm.CanManageChat {
				result = "Да"
			}
			return result, chm.CanManageChat
		}},
		Criterium{"Отправка сообщений", func(chm tg.ChatMember) (string, bool) {
			result := "Нет"
			if chm.CanPostMessages {
				result = "Да"
			}
			return result, chm.CanPostMessages
		}},
		Criterium{"Изменение сообщений", func(chm tg.ChatMember) (string, bool) {
			result := "Нет"
			if chm.CanEditMessages {
				result = "Да"
			}
			return result, chm.CanEditMessages
		}},
		Criterium{"Удаление сообщений", func(chm tg.ChatMember) (string, bool) {
			result := "Нет"
			if chm.CanDeleteMessages {
				result = "Да"
			}
			return result, chm.CanDeleteMessages
		}},
	}

	data := struct {
		AccessRequirement string
		TgPublicId        int64
	}{
		AccessRequirement: generateAccessRequirement(newMember, criteria),
		TgPublicId:        chatMemberUpdated.Chat.ID,
	}

	h.SendTemplate(chatMemberUpdated.Chat.ID, "msg-add-bot-to-channel.html", data)
}
