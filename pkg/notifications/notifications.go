package notifications

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/slntopp/core-chatting/cc"
	"github.com/slntopp/core-chatting/pkg/notifications/gateways"
	"github.com/slntopp/core-chatting/pkg/notifications/types"
	"go.uber.org/zap"
	"strconv"
)

type NotificationDispatcher interface {
	Notify(event string, payload interface{})
}

var events = map[cc.EventType]string{
	cc.EventType_CHAT_CREATED: "chat_created",
}

func EventToNotification(eventType cc.EventType) string {
	return events[eventType]
}

type NotificationGateway interface {
	Send(event string, meta types.EventMeta, payload interface{}) error
}

type GatewayConfig struct {
	Type        string                     `yaml:"type" json:"type"`
	Events      map[string]types.EventMeta `yaml:"events" json:"events"`
	Credentials map[string]string          `yaml:"credentials" json:"credentials"`
}

type GatewayHolder struct {
	Gateway NotificationGateway
	Events  map[string]types.EventMeta
}

type notificationDispatcher struct {
	log      *zap.Logger
	Gateways []GatewayHolder
}

func NewNotificationDispatcher(_log *zap.Logger, configs []GatewayConfig) NotificationDispatcher {
	log := _log.Named("NotificationDispatcher")

	var gws []GatewayHolder

	for _, cfg := range configs {
		switch cfg.Type {
		case "telegram":
			token, ok := cfg.Credentials["token"]
			if !ok {
				log.Fatal("No telegram token in credentials")
			}
			chatID, err := strconv.Atoi(cfg.Credentials["chat_id"])
			if err != nil {
				log.Fatal("No valid telegram chat id in credentials")
			}
			bot, err := tgbotapi.NewBotAPI(token)
			if err != nil {
				log.Fatal("Failed to initialize telegram bot", zap.Error(err))
			}
			tg := gateways.NewTelegramGateway(log, bot, int64(chatID))
			gws = append(gws, GatewayHolder{
				Gateway: tg,
				Events:  cfg.Events,
			})
		default:
			log.Error("Unknown gateway type. Skip...", zap.String("type", cfg.Type))
			continue
		}
		log.Debug("Gateway initialized", zap.String("type", cfg.Type), zap.Any("events", cfg.Events))
	}
	return &notificationDispatcher{log: log, Gateways: gws}
}

func (nd *notificationDispatcher) Notify(event string, payload interface{}) {
	for _, gh := range nd.Gateways {
		for key, meta := range gh.Events {
			if key == event {
				if err := gh.Gateway.Send(event, meta, payload); err != nil {
					nd.log.Error("Failed to send message to gateway", zap.Error(err))
				}
				break
			}
		}
	}
}
