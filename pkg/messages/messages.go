package messages

import (
	"context"
	"errors"
	"fmt"
	structpb "google.golang.org/protobuf/types/known/structpb"
	"time"

	"github.com/slntopp/core-chatting/pkg/pubsub"

	"github.com/slntopp/core-chatting/cc"
	"github.com/slntopp/core-chatting/pkg/core"
	"github.com/slntopp/core-chatting/pkg/graph"

	"connectrpc.com/connect"
	"go.uber.org/zap"
)

type MessagesServer struct {
	log *zap.Logger

	chatCtrl *graph.ChatsController
	msgCtrl  *graph.MessagesController
	ps       *pubsub.PubSub
}

func NewMessagesServer(logger *zap.Logger, chatCtrl *graph.ChatsController, msgCtrl *graph.MessagesController, ps *pubsub.PubSub) *MessagesServer {
	return &MessagesServer{log: logger.Named("MessagesServer"), chatCtrl: chatCtrl, msgCtrl: msgCtrl, ps: ps}
}

func (s *MessagesServer) Get(ctx context.Context, req *connect.Request[cc.Chat]) (*connect.Response[cc.Messages], error) {
	log := s.log.Named("Get")
	log.Debug("Request received", zap.Any("req", req.Msg))

	requestor := ctx.Value(core.ChatAccount).(string)

	chat, err := s.chatCtrl.Get(ctx, req.Msg.GetUuid(), requestor)
	if err != nil {
		return nil, err
	}

	if chat.Role < cc.Role_USER {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("no access to chat"))
	}

	messages, err := s.chatCtrl.GetMessages(ctx, req.Msg, chat.Role == cc.Role_ADMIN)
	if err != nil {
		return nil, err
	}

	for index := range messages {
		if !core.In(requestor, messages[index].Readers) {
			newMessage, err := s.msgCtrl.Read(ctx, messages[index], requestor)
			if err != nil {
				log.Error("Failed to update reader", zap.Error(err))
			}
			messages[index] = newMessage
		}
	}

	resp := connect.NewResponse[cc.Messages](&cc.Messages{
		Messages: messages,
	})

	go s.ps.Pub(ctx, requestor, &cc.Event{Type: cc.EventType_CHAT_READ, Item: &cc.Event_Chat{Chat: req.Msg}})

	return resp, nil
}

func (s *MessagesServer) Send(ctx context.Context, req *connect.Request[cc.Message]) (*connect.Response[cc.Message], error) {
	log := s.log.Named("Send")
	log.Debug("Request received", zap.Any("req", req.Msg))

	requestor := ctx.Value(core.ChatAccount).(string)

	msg := req.Msg

	chat, err := s.chatCtrl.Get(ctx, msg.GetChat(), requestor)
	if err != nil {
		return nil, err
	}
	log.Debug("Chat retrieved", zap.String("requestor", requestor), zap.String("role", chat.Role.String()))

	if chat.Role < cc.Role_USER {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("no access to chat"))
	}

	msg.Sender = requestor

	if (msg.Kind == cc.Kind_ADMIN_ONLY || msg.UnderReview) && chat.Role != cc.Role_ADMIN {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("can't send admin only message"))
	}

	msg.Sent = time.Now().UnixMilli()

	if chat.GetGateways() != nil {
		if msg.GetMeta() == nil {
			msg.Meta = map[string]*structpb.Value{}
		}
		data := chat.GetMeta().GetData()
		if data != nil {
			for _, gate := range chat.GetGateways() {
				msg.Meta[fmt.Sprintf("%s_chat_id", gate)] = data[gate]
			}
		}
	}

	message, err := s.msgCtrl.Send(ctx, msg)
	if err != nil {
		return nil, err
	}

	go handleNotify(ctx, log, s.ps, message, chat, cc.EventType_MESSAGE_SENT)

	resp := connect.NewResponse[cc.Message](message)

	return resp, nil
}

func (s *MessagesServer) Update(ctx context.Context, req *connect.Request[cc.Message]) (*connect.Response[cc.Message], error) {
	log := s.log.Named("Update")
	log.Debug("Request received", zap.Any("req", req.Msg))

	requestor := ctx.Value(core.ChatAccount).(string)

	chat, err := s.chatCtrl.Get(ctx, req.Msg.GetChat(), requestor)
	if err != nil {
		return nil, err
	}

	if chat.Role < cc.Role_USER {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("no access to chat"))
	}

	if req.Msg.Kind == cc.Kind_ADMIN_ONLY && chat.Role != cc.Role_ADMIN {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("can't send admin only message"))
	}

	if requestor != req.Msg.GetSender() && !core.In(requestor, chat.GetAdmins()) {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("you are not message sender or chat admin"))
	}

	req.Msg.Edited = time.Now().UnixMilli()

	oldMessage, err := s.msgCtrl.Get(ctx, req.Msg.GetUuid())

	if err != nil {
		return nil, err
	}

	message, err := s.msgCtrl.Update(ctx, req.Msg)
	if err != nil {
		return nil, err
	}

	if oldMessage.GetKind() != message.GetKind() || oldMessage.GetUnderReview() != message.GetUnderReview() {
		go handleSpecialNotify(ctx, log, s.ps, message, oldMessage, chat)
	} else {
		go handleNotify(ctx, log, s.ps, message, chat, cc.EventType_MESSAGE_UPDATED)
	}

	resp := connect.NewResponse[cc.Message](message)

	return resp, nil
}

func (s *MessagesServer) Delete(ctx context.Context, req *connect.Request[cc.Message]) (*connect.Response[cc.Message], error) {
	log := s.log.Named("Delete")
	log.Debug("Request received", zap.Any("req", req.Msg))

	requestor := ctx.Value(core.ChatAccount).(string)

	chat, err := s.chatCtrl.Get(ctx, req.Msg.GetChat(), requestor)
	if err != nil {
		return nil, err
	}

	if chat.Role < cc.Role_USER {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("no access to chat"))
	}

	if requestor != req.Msg.GetSender() && !core.In(requestor, chat.GetAdmins()) {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("you are not message sender or chat admin"))
	}

	message, err := s.msgCtrl.Delete(ctx, req.Msg)
	if err != nil {
		return nil, err
	}

	go handleNotify(ctx, log, s.ps, message, chat, cc.EventType_MESSAGE_DELETED)

	resp := connect.NewResponse[cc.Message](message)

	return resp, nil
}

func (s *MessagesServer) GetByGateway(ctx context.Context, req *connect.Request[cc.GetawayRequest]) (*connect.Response[cc.Message], error) {
	log := s.log.Named("GetByGateway")
	log.Debug("Request received", zap.Any("req", req.Msg))

	requestor := ctx.Value(core.ChatAccount).(string)

	message, err := s.msgCtrl.GetByGateway(ctx, req.Msg)
	if err != nil {
		return nil, err
	}

	if message.Sender != requestor {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("no access to message"))
	}

	resp := connect.NewResponse[cc.Message](message)

	return resp, nil
}

func handleSpecialNotify(ctx context.Context, log *zap.Logger, ps *pubsub.PubSub, msg *cc.Message, oldMsg *cc.Message, chat *cc.Chat) {

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

func handleNotify(ctx context.Context, log *zap.Logger, ps *pubsub.PubSub, msg *cc.Message, chat *cc.Chat, eventType cc.EventType) {
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
