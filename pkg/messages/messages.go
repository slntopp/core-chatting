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

	if !In(requestor, req.Msg.GetUsers()) {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("no access to chat"))
	}

	messages, err := s.chatCtrl.GetMessages(ctx, req.Msg)
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
	}

	if !In(requestor, chat.GetUsers()) {
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
