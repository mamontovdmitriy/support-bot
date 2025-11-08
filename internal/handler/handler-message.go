package handler

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	reExtractUserId = regexp.MustCompile(`(?i)User\s?ID:\s?([-\d]+)`)
)

/**
 * Пересылка сообщения клиента из бота в комментарий соответствующего поста в группе обсуждения
 * Оратная пересылка сообщения менеджера из комментов соответствующему клиенту.
 */
func (c *Handler) proxyMessage(message *tg.Message) {
	// Create forward post into channel - save forward post id
	if message.ForwardFromChat != nil && message.ForwardFromChat.ID == c.cfg.PublicId {
		userId, err := extractUserId(message.Text)
		if err != nil {
			c.services.Log.Warn("Handler.proxyMessage: userId not found")
			return
		}

		err = c.services.UserInfoPost.SaveUserInfoPost(userId, int64(message.MessageID))
		if err != nil {
			c.services.Log.Errorf("Handler.proxyMessage: error saving forward post ID %v", err)
		}
		// update async message.MessageID
		return
	}

	if message.From == nil || message.From.ID == c.cfg.SystemUserId {
		c.services.Log.Info("Handler.proxyMessage: ignore messages from Telegram")
		return // ignore messages from Telegram
	}

	// response to bot
	if message.ReplyToMessage != nil {
		userId, err := extractUserId(message.ReplyToMessage.Text)
		if err != nil {
			c.services.Log.Warn("Handler.proxyMessage: userId not found in ReplyToMessage")
			return
		}
		_, err = c.bot.CopyMessage(tg.NewCopyMessage(userId, message.Chat.ID, message.MessageID))
		if err != nil {
			c.services.Log.Errorf("Handler.proxyMessage: error coping reply message to bot - %v", err)
			// show resend button
		}
		return
	}

	// ответ в группе без указания адресата
	if message.SenderChat != nil {
		msgText := fmt.Sprintf("Сообщение не доставлено!\n\nНе указат получатель ответа.\nДля этого отвечайте на сообщение с 'User ID: ***'.")
		_, err := c.bot.Send(tg.NewMessage(c.cfg.ChannelId, msgText))
		if err != nil {
			c.services.Log.Warnf("Handler.proxyMessage: error sending hint  %v", err)
			// show resend button
			return
		}
		c.services.Log.Info("Handler.proxyMessage: ignore messages from Telegram")
		return // ignore messages from Telegram
	}

	// Send message from bot to channel
	userId := message.From.ID
	forwardPostId, err := c.services.UserInfoPost.GetForwardId(userId)
	if err != nil {
		// Create user info post
		msgText := fmt.Sprintf("User ID: %d\n--\n", userId)
		_, err := c.bot.Send(tg.NewMessage(c.cfg.PublicId, msgText))
		if err != nil {
			c.services.Log.Warnf("Handler.proxyMessage: error send post %v", err)
			// show resend button
			return
		}

		// Wait for creation of forward post
		go func() {
			for {
				forwardPostId, err := c.services.UserInfoPost.GetForwardId(userId)
				if err == nil {
					// copy message from Bot into forward post comments
					c.copyMessage(forwardPostId, message)
					c.services.Log.Infof("Message async proxied %d", message.MessageID)
					return
				}
				time.Sleep(3 * time.Second)
			}
		}()
	} else {
		// copy message from Bot into forward post comments
		c.copyMessage(forwardPostId, message)
		c.services.Log.Infof("Message proxied %d", message.MessageID)
	}
}

func (c *Handler) copyMessage(forwardPostId int64, message *tg.Message) {
	_, err := c.bot.CopyMessage(tg.CopyMessageConfig{
		BaseChat: tg.BaseChat{
			ChatID:           c.cfg.ChannelId,
			ReplyToMessageID: int(forwardPostId),
		},
		FromChatID: message.Chat.ID,
		MessageID:  message.MessageID,
	})
	if err != nil {
		c.services.Log.Errorf("Handler.proxyMessage: error coping message to chat - %v", err)
		// show resend button
	}
}

func extractUserId(text string) (int64, error) {
	match := reExtractUserId.FindStringSubmatch(text)

	if len(match) == 0 {
		return 0, errors.New("regexp failed - userId not found")
	}

	return strconv.ParseInt(match[1], 10, 64)
}
