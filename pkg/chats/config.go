package chats

import (
	"context"
	"github.com/slntopp/nocloud-proto/access"
	settingspb "github.com/slntopp/nocloud-proto/settings"
	sc "github.com/slntopp/nocloud/pkg/settings/client"
	"go.uber.org/zap"
)

const (
	ticketsKey string = "tickets-settings"
)

type TicketsSettingsConf struct {
	RoutineFrequency             int `json:"routine-frequency"`                // In seconds
	CloseInactiveChatsAfterHours int `json:"close-inactive-chats-after-hours"` // In hours
}

var (
	ticketsSetting = &sc.Setting[TicketsSettingsConf]{
		Value: TicketsSettingsConf{
			RoutineFrequency:             3600,
			CloseInactiveChatsAfterHours: 120,
		},
		Description: "Tickets Settings",
		Level:       access.Level_ADMIN,
	}
)

func MakeTicketsConf(ctx context.Context, log *zap.Logger, settingsClient *settingspb.SettingsServiceClient) (conf TicketsSettingsConf) {
	sc.Setup(log, ctx, settingsClient)

	if err := sc.Fetch(ticketsKey, &conf, ticketsSetting); err != nil {
		log.Warn("Failed to fetch tickets settings. Setting default", zap.Error(err))
		conf = ticketsSetting.Value
	}

	log.Debug("Got tickets config", zap.Any("conf", conf))
	return conf
}
