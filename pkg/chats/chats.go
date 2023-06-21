package chats

import (
	"context"

	"github.com/slntopp/core-chatting/cc"
	"github.com/slntopp/core-chatting/pkg/graph"

	"github.com/bufbuild/connect-go"
	"go.uber.org/zap"
)

type ChatsServer struct {
	log *zap.Logger

	ctrl *graph.ChatsController
}

func NewChatsServer(logger *zap.Logger, ctrl *graph.ChatsController) *ChatsServer {
	return &ChatsServer{log: logger.Named("ChatsServer"), ctrl: ctrl}
}

func (s *ChatsServer) Create(ctx context.Context, req *connect.Request[cc.Chat]) (*connect.Response[cc.Chat], error) {
	log := s.log.Named("Create")
	log.Debug("Request received", zap.Any("req", req.Msg))

	chat, err := s.ctrl.Create(ctx, req.Msg)
	if err != nil {
		return nil, err
	}

	resp := connect.NewResponse[cc.Chat](chat)

	return resp, nil
}

func (s *ChatsServer) Update(ctx context.Context, req *connect.Request[cc.Chat]) (*connect.Response[cc.Chat], error) {
	log := s.log.Named("Create")
	log.Debug("Request received", zap.Any("req", req.Msg))

	chat, err := s.ctrl.Update(ctx, req.Msg)
	if err != nil {
		return nil, err
	}

	resp := connect.NewResponse[cc.Chat](chat)

	return resp, nil
}

func (s *ChatsServer) List(ctx context.Context, req *connect.Request[cc.Empty]) (*connect.Response[cc.Chats], error) {
	log := s.log.Named("Create")
	log.Debug("Request received", zap.Any("req", req.Msg))

	chats, err := s.ctrl.List(ctx)
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

	chat, err := s.ctrl.Delete(ctx, req.Msg)
	if err != nil {
		return nil, err
	}

	resp := connect.NewResponse[cc.Chat](chat)

	return resp, nil
}
