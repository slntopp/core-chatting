package pubsub

import (
	"context"
	"github.com/slntopp/core-chatting/cc"
	"go.uber.org/zap"
)

func HandleSpecialNotify(ctx context.Context, log *zap.Logger, ps *PubSub, msg *cc.Message, oldMsg *cc.Message, chat *cc.Chat) {

	var adminEvent = &cc.Event{
		Type: cc.EventType_MESSAGE_UPDATED,
		Item: &cc.Event_Msg{Msg: msg},
	}

	var userEvent = &cc.Event{
		Item: &cc.Event_Msg{Msg: msg},
	}

	for _, admin := range chat.GetAdmins() {
		go ps.Pub(ctx, admin, adminEvent)
	}

	go ps.PubGateway(ctx, adminEvent, chat.GetGateways())

	diffKinds := msg.GetKind() != oldMsg.GetKind()
	diffReviews := msg.GetUnderReview() != oldMsg.GetUnderReview()

	sendEvent := (diffKinds || diffReviews) && msg.GetKind() == cc.Kind_DEFAULT && !msg.GetUnderReview()

	deleteEvent := (diffKinds && msg.GetKind() == cc.Kind_ADMIN_ONLY) ||
		(diffReviews && msg.GetUnderReview())

	skip := msg.GetKind() == cc.Kind_DEFAULT && oldMsg.GetKind() == cc.Kind_ADMIN_ONLY && msg.GetUnderReview() && !oldMsg.GetUnderReview()

	if sendEvent {
		userEvent.Type = cc.EventType_MESSAGE_SENT
	} else if deleteEvent && !skip {
		userEvent.Type = cc.EventType_MESSAGE_DELETED
	} else {
		return
	}

	for _, user := range chat.GetUsers() {
		go ps.Pub(ctx, user, userEvent)
	}
}

func HandleNotifyMessage(ctx context.Context, log *zap.Logger, ps *PubSub, msg *cc.Message, chat *cc.Chat, eventType cc.EventType) {
	var event = &cc.Event{
		Type: eventType,
		Item: &cc.Event_Msg{Msg: msg},
	}

	if msg.Kind == cc.Kind_DEFAULT && !msg.UnderReview {
		for _, user := range chat.GetUsers() {
			go ps.Pub(ctx, user, event)
		}
	}

	for _, admin := range chat.GetAdmins() {
		go ps.Pub(ctx, admin, event)
	}

	go ps.PubGateway(ctx, event, chat.GetGateways())
}

func HandleNotifyChat(ctx context.Context, log *zap.Logger, ps *PubSub, chat *cc.Chat, eventType cc.EventType) {
	var event = &cc.Event{
		Type: eventType,
		Item: &cc.Event_Chat{Chat: chat},
	}

	for _, user := range chat.GetUsers() {
		log.Info("Send to", zap.Any("User", user))
		go ps.Pub(ctx, user, event)
	}

	for _, admin := range chat.GetAdmins() {
		log.Info("Send to", zap.Any("Admin", admin))
		go ps.Pub(ctx, admin, event)
	}

	go ps.PubGateway(ctx, event, chat.GetGateways())
}
