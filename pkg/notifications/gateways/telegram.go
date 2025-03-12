package gateways

import (
	"fmt"
	"github.com/slntopp/core-chatting/cc"
	"github.com/slntopp/core-chatting/pkg/notifications/params"
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

func matchingDepartments(chat *cc.Chat, meta types.EventMeta) bool {
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
			return false
		}
	}
	return true
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
		if !matchingDepartments(chat, meta) {
			log.Debug("Departments not matched. Skip...")
			return nil
		}
		tgChat := &destination{id: fmt.Sprintf("-%d", tg.chatID)}
		result := params.ParseParameters(meta.Message, params.ParametersFromChat(chat))
		log.Debug("result answer", zap.String("answer", result))
		_, err := tg.bot.Send(tgChat, result, tele.ModeMarkdown)
		if err != nil {
			return fmt.Errorf("failed to send telegram message: %w", err)
		}

	case "sla_violation_level_1":
		chat, ok := payload.(*cc.Chat)
		if !ok {
			return fmt.Errorf("failed to obtain created chat from payload")
		}
		if !matchingDepartments(chat, meta) {
			log.Debug("Departments not matched. Skip...")
			return nil
		}
		tgChat := &destination{id: fmt.Sprintf("-%d", tg.chatID)}
		result := params.ParseParameters(meta.Message, params.ParametersFromChat(chat))
		log.Debug("result answer", zap.String("answer", result))
		_, err := tg.bot.Send(tgChat, result, tele.ModeMarkdown)
		if err != nil {
			return fmt.Errorf("failed to send telegram message: %w", err)
		}

	default:
		log.Warn("Cannot handle this event")
	}

	return nil
}

type destination struct {
	id string
}

func (d destination) Recipient() string {
	return d.id
}
