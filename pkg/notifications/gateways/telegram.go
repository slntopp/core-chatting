package gateways

import (
	"fmt"
	"github.com/slntopp/core-chatting/cc"
	"github.com/slntopp/core-chatting/pkg/notifications/types"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v4"
	"strings"
)

type TelegramGateway struct {
	log    *zap.Logger
	bot    *tele.Bot
	chatID int64
}

func NewTelegramGateway(log *zap.Logger, bot *tele.Bot, chatID int64) *TelegramGateway {
	return &TelegramGateway{log: log.Named("TelegramGateway"), bot: bot, chatID: chatID}
}

func (tg *TelegramGateway) Send(event string, meta types.EventMeta, payload interface{}) error {
	log := tg.log.With(zap.String("event", event))
	log.Debug("New event", zap.Any("meta", meta), zap.Any("payload", payload))

	switch event {
	case "chat_created":
		chat, ok := payload.(*cc.Chat)
		if !ok {
			return fmt.Errorf("failed to obtain created chat from payload")
		}
		depStr, _ := meta.Meta["departments"].(string)
		deps := strings.Split(depStr, ",")
		if len(deps) > 0 && depStr != "" {
			found := false
			for _, dep := range deps {
				if dep == chat.Department {
					found = true
					break
				}
			}
			if !found {
				log.Debug("Department not matched. Skip...")
				return nil
			}
		}
		tgChat := &tele.Chat{ID: tg.chatID}
		_, err := tg.bot.Send(tgChat, meta.Message)
		if err != nil {
			return fmt.Errorf("failed to send telegram message: %w", err)
		}

	default:
		log.Warn("Cannot handle this event")
	}

	return nil
}
