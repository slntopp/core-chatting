package chats

import (
	"context"
	"errors"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/slntopp/core-chatting/pkg/core"
	"github.com/slntopp/core-chatting/pkg/pubsub"

	"github.com/slntopp/core-chatting/cc"
	"github.com/slntopp/core-chatting/pkg/graph"

	"connectrpc.com/connect"
	"go.uber.org/zap"
)

type ChatsServer struct {
	log *zap.Logger

	ctrl       *graph.ChatsController
	users_ctrl *graph.UsersController
	ps         *pubsub.PubSub
}

func NewChatsServer(logger *zap.Logger, ctrl *graph.ChatsController, users_ctrl *graph.UsersController, ps *pubsub.PubSub) *ChatsServer {
	return &ChatsServer{log: logger.Named("ChatsServer"), ctrl: ctrl, users_ctrl: users_ctrl, ps: ps}
}

func (s *ChatsServer) Create(ctx context.Context, req *connect.Request[cc.Chat]) (*connect.Response[cc.Chat], error) {
	log := s.log.Named("Create")
	log.Debug("Request received", zap.Any("req", req.Msg))

	requestor := ctx.Value(core.ChatAccount).(string)

	req.Msg.Owner = requestor

	msg := req.Msg
	resolve, err := s.users_ctrl.Resolve(ctx, []string{requestor})
	if err != nil {
		return nil, err
	}
	user := resolve[0]

	if user.GetData() != nil {
		fields := user.GetData().GetFields()

		log.Debug("Field", zap.Any("f", fields))

		if msg.GetMeta() == nil {
			msg.Meta = &cc.ChatMeta{
				Data: map[string]*structpb.Value{},
			}
		}

		if msg.GetMeta().GetData() == nil {
			msg.Meta.Data = map[string]*structpb.Value{}
		}

		if fields != nil {
			for _, gate := range msg.GetGateways() {
				log.Debug("gate", zap.String("g", gate))
				if val, ok := fields[gate]; ok {
					log.Debug("ok")
					msg.Meta.Data[gate] = val
				}
			}
		}
		err := s.ctrl.DeleteGateways(ctx, fields)
		if err != nil {
			log.Error("Failed to delete old gates", zap.Error(err))
			return nil, err
		}
	}

	chat, err := s.ctrl.Create(ctx, req.Msg)
	if err != nil {
		return nil, err
	}

	go pubsub.HandleNotifyTicket(ctx, log, s.ps, chat)
	go pubsub.HandleNotifyChat(ctx, log, s.ps, chat, cc.EventType_CHAT_CREATED)

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

	go pubsub.HandleNotifyChat(ctx, log, s.ps, chat, cc.EventType_CHAT_UPDATED)

	resp := connect.NewResponse[cc.Chat](chat)

	return resp, nil
}

func (s *ChatsServer) Get(ctx context.Context, req *connect.Request[cc.Chat]) (*connect.Response[cc.Chat], error) {
	log := s.log.Named("Create")
	log.Debug("Request received", zap.Any("req", req.Msg))

	requestor := ctx.Value(core.ChatAccount).(string)

	chat, err := s.ctrl.Get(ctx, req.Msg.Uuid, requestor)
	if err != nil {
		return nil, err
	}

	if chat.Role == cc.Role_NOACCESS {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("no access to chat"))
	}

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

	go pubsub.HandleNotifyChat(ctx, log, s.ps, chat, cc.EventType_CHAT_DELETED)

	resp := connect.NewResponse[cc.Chat](chat)

	return resp, nil
}
