package pubsub

import (
	"context"
	"github.com/slntopp/core-chatting/cc"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
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

	go ps.PubGateway(ctx, adminEvent, append(chat.GetGateways(), "whmcs"))

	diffKinds := msg.GetKind() != oldMsg.GetKind()
	diffReviews := msg.GetUnderReview() != oldMsg.GetUnderReview()

	sendEvent := (diffKinds || diffReviews) && msg.GetKind() == cc.Kind_DEFAULT && !msg.GetUnderReview()

	if diffReviews && !msg.GetUnderReview() {
		if msg.Meta == nil {
			msg.Meta = make(map[string]*structpb.Value)
		}
		msg.Meta["approved"] = structpb.NewBoolValue(true)
	}

	// A message that turned operator-only was visible to the customer a moment
	// ago, so their client has to be told to drop it - the copilot lane needs this
	// exactly as much as a note does.
	deleteEvent := (diffKinds && cc.IsOperatorOnly(msg.GetKind())) ||
		(diffReviews && msg.GetUnderReview())

	skip := msg.GetKind() == cc.Kind_DEFAULT && cc.IsOperatorOnly(oldMsg.GetKind()) && msg.GetUnderReview() && !oldMsg.GetUnderReview()

	if sendEvent {
		userEvent.Type = cc.EventType_MESSAGE_SENT
	} else if deleteEvent && !skip {
		userEvent.Type = cc.EventType_MESSAGE_DELETED
		// The customer only needs to know which message to drop, not what the
		// note now says.
		userEvent.Item = &cc.Event_Msg{Msg: &cc.Message{Uuid: msg.GetUuid(), Chat: msg.Chat}}
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

	go ps.PubGateway(ctx, event, append(chat.GetGateways(), "whmcs"))
}

func HandleNotifyChat(ctx context.Context, log *zap.Logger, ps *PubSub, chat *cc.Chat, eventType cc.EventType) {
	var event = &cc.Event{
		Type: eventType,
		Item: &cc.Event_Chat{Chat: chat},
	}

	userEvent := &cc.Event{
		Type: eventType,
		Item: &cc.Event_Chat{Chat: chatForUsers(chat)},
	}

	for _, user := range chat.GetUsers() {
		log.Info("Send to", zap.Any("User", user))
		go ps.Pub(ctx, user, userEvent)
	}

	for _, admin := range chat.GetAdmins() {
		log.Info("Send to", zap.Any("Admin", admin))
		go ps.Pub(ctx, admin, event)
	}

	go ps.PubGateway(ctx, event, append(chat.GetGateways(), "whmcs"))
}

// chatForUsers strips the first/last message off a chat headed to customers when
// it is one they must not see. The chat here is whatever the caller holds -
// often the object an operator's client sent to Update/ChangeStatus, whose meta
// was computed for the operator and so may carry an internal note.
func chatForUsers(chat *cc.Chat) *cc.Chat {
	meta := chat.GetMeta()
	if !isHidden(meta.GetLastMessage()) && !isHidden(meta.GetFirstMessage()) {
		return chat
	}

	clone := proto.Clone(chat).(*cc.Chat)
	if isHidden(clone.Meta.LastMessage) {
		clone.Meta.LastMessage = nil
	}
	if isHidden(clone.Meta.FirstMessage) {
		clone.Meta.FirstMessage = nil
	}
	return clone
}

func isHidden(msg *cc.Message) bool {
	return msg != nil && (cc.IsOperatorOnly(msg.GetKind()) || msg.GetUnderReview())
}
