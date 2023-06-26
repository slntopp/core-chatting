package chats

import (
	"context"
	"errors"
	"github.com/slntopp/core-chatting/pkg/core"
	"github.com/slntopp/core-chatting/pkg/pubsub"

	"github.com/slntopp/core-chatting/cc"
	"github.com/slntopp/core-chatting/pkg/graph"

	"github.com/bufbuild/connect-go"
	"go.uber.org/zap"
)

type ChatsServer struct {
	log *zap.Logger

	ctrl *graph.ChatsController
	ps   *pubsub.PubSub
}

func NewChatsServer(logger *zap.Logger, ctrl *graph.ChatsController, ps *pubsub.PubSub) *ChatsServer {
	return &ChatsServer{log: logger.Named("ChatsServer"), ctrl: ctrl, ps: ps}
}

func (s *ChatsServer) Create(ctx context.Context, req *connect.Request[cc.Chat]) (*connect.Response[cc.Chat], error) {
	log := s.log.Named("Create")
	log.Debug("Request received", zap.Any("req", req.Msg))

	requestor := ctx.Value(core.ChatAccount).(string)

	req.Msg.Owner = requestor

	chat, err := s.ctrl.Create(ctx, req.Msg)
	if err != nil {
		return nil, err
	}

	go handleNotify(ctx, s.ps, chat, cc.EventType_CHAT_CREATED)

	resp := connect.NewResponse[cc.Chat](chat)

	return resp, nil
}

func (s *ChatsServer) Update(ctx context.Context, req *connect.Request[cc.Chat]) (*connect.Response[cc.Chat], error) {
	log := s.log.Named("Create")
	log.Debug("Request received", zap.Any("req", req.Msg))

	requestor := ctx.Value(core.ChatAccount).(string)

	chat, err := s.ctrl.Get(ctx, req.Msg.Uuid, requestor)
	if err != nil {
		return nil, err
	}

	if chat.Role < cc.Role_OWNER {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("no access to chat"))
	}

	chat, err = s.ctrl.Update(ctx, req.Msg)
	if err != nil {
		return nil, err
	}

	go handleNotify(ctx, s.ps, chat, cc.EventType_CHAT_UPDATED)

	resp := connect.NewResponse[cc.Chat](chat)

	return resp, nil
}

func (s *ChatsServer) List(ctx context.Context, req *connect.Request[cc.Empty]) (*connect.Response[cc.Chats], error) {
	log := s.log.Named("Create")
	log.Debug("Request received", zap.Any("req", req.Msg))

	requestor := ctx.Value(core.ChatAccount).(string)

	chats, err := s.ctrl.List(ctx, requestor)
	if err != nil {
		return nil, err
	}

	resp := connect.NewResponse[cc.Chats](&cc.Chats{
		Chats: chats,
	})

	return resp, nil
}

func (s *ChatsServer) Delete(ctx context.Context, req *connect.Request[cc.Chat]) (*connect.Response[cc.Chat], error) {
	log := s.log.Named("Create")
	log.Debug("Request received", zap.Any("req", req.Msg))

	requestor := ctx.Value(core.ChatAccount).(string)

	chat, err := s.ctrl.Get(ctx, req.Msg.Uuid, requestor)
	if err != nil {
		return nil, err
	}

	if chat.Role < cc.Role_OWNER {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("no access to chat"))
	}

	chat, err = s.ctrl.Delete(ctx, req.Msg)
	if err != nil {
		return nil, err
	}

	go handleNotify(ctx, s.ps, chat, cc.EventType_CHAT_DELETED)

	resp := connect.NewResponse[cc.Chat](chat)

	return resp, nil
}

func handleNotify(ctx context.Context, ps *pubsub.PubSub, chat *cc.Chat, eventType cc.EventType) {
	var event = &cc.Event{
		Type: eventType,
		Item: &cc.Event_Chat{Chat: chat},
	}

	for _, user := range chat.Users {
		go ps.Pub(ctx, user, event)
	}

	for _, admin := range chat.Admins {
		go ps.Pub(ctx, admin, event)
	}

}
