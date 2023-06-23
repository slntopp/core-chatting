package messages

import (
	"context"
	"errors"

	"github.com/slntopp/core-chatting/cc"
	"github.com/slntopp/core-chatting/pkg/core"
	"github.com/slntopp/core-chatting/pkg/graph"

	"github.com/bufbuild/connect-go"
	"go.uber.org/zap"
)

type MessagesServer struct {
	log *zap.Logger

	chatCtrl *graph.ChatsController
	msgCtrl  *graph.MessagesController
}

func NewMessagesServer(logger *zap.Logger, chatCtrl *graph.ChatsController, msgCtrl *graph.MessagesController) *MessagesServer {
	return &MessagesServer{log: logger.Named("MessagesServer"), chatCtrl: chatCtrl, msgCtrl: msgCtrl}
}

func In[T comparable](o T, a []T) bool {
	for _, el := range a {
		if el == o {
			return true
		}
	}
	return false
}

func (s *MessagesServer) Get(ctx context.Context, req *connect.Request[cc.Chat]) (*connect.Response[cc.Messages], error) {
	log := s.log.Named("Get")
	log.Debug("Request received", zap.Any("req", req.Msg))

	requestor := ctx.Value(core.ChatAccount).(string)

	chat, err := s.chatCtrl.Get(ctx, req.Msg.GetUuid())
	if err != nil {
		return nil, err
	}

	is_user := In(requestor, chat.GetUsers())
	is_admin := In(requestor, chat.GetAdmins())

	if !is_user && !is_admin {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("no access to chat"))
	}

	messages, err := s.chatCtrl.GetMessages(ctx, req.Msg, is_admin)
	if err != nil {
		return nil, err
	}

	resp := connect.NewResponse[cc.Messages](&cc.Messages{
		Messages: messages,
	})

	return resp, nil
}

func (s *MessagesServer) Send(ctx context.Context, req *connect.Request[cc.Message]) (*connect.Response[cc.Message], error) {
	log := s.log.Named("Send")
	log.Debug("Request received", zap.Any("req", req.Msg))

	requestor := ctx.Value(core.ChatAccount).(string)

	chat, err := s.chatCtrl.Get(ctx, req.Msg.GetChat())

	if err != nil {
		return nil, err
	}

	if requestor != req.Msg.GetSender() {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("no access to chat"))

	is_user := In(requestor, chat.GetUsers())
	is_admin := In(requestor, chat.GetAdmins())

	if req.Msg.Kind == cc.Kind_ADMIN_ONLY && !is_admin {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("can't send admin only message"))
	}

	if !is_user && !is_admin {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("no access to chat"))
	}

	message, err := s.msgCtrl.Send(ctx, req.Msg)
	if err != nil {
		return nil, err
	}

	resp := connect.NewResponse[cc.Message](message)

	return resp, nil
}

func (s *MessagesServer) Update(ctx context.Context, req *connect.Request[cc.Message]) (*connect.Response[cc.Message], error) {
	log := s.log.Named("Update")
	log.Debug("Request received", zap.Any("req", req.Msg))

	requestor := ctx.Value(core.ChatAccount).(string)

	if requestor != req.Msg.GetSender() {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("no access to chat"))
	}

	message, err := s.msgCtrl.Update(ctx, req.Msg)
	if err != nil {
		return nil, err
	}

	resp := connect.NewResponse[cc.Message](message)

	return resp, nil
}

func (s *MessagesServer) Delete(ctx context.Context, req *connect.Request[cc.Message]) (*connect.Response[cc.Message], error) {
	log := s.log.Named("Delete")
	log.Debug("Request received", zap.Any("req", req.Msg))

	requestor := ctx.Value(core.ChatAccount).(string)

	chat, err := s.chatCtrl.Get(ctx, req.Msg.GetChat())

	if err != nil {
		return nil, err
	}

	if requestor != req.Msg.GetSender() {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("no access to chat"))
	}

	if !In(requestor, chat.GetUsers()) {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("no access to chat"))
	}

	message, err := s.msgCtrl.Delete(ctx, req.Msg)
	if err != nil {
		return nil, err
	}

	resp := connect.NewResponse[cc.Message](message)

	return resp, nil
}
