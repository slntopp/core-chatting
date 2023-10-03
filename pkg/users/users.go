package users

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/slntopp/core-chatting/cc"
	"github.com/slntopp/core-chatting/pkg/core"
	"github.com/slntopp/core-chatting/pkg/graph"
	"go.uber.org/zap"
)

type UsersServer struct {
	log *zap.Logger

	ctrl *graph.UsersController
}

func NewUsersServer(logger *zap.Logger, ctrl *graph.UsersController) *UsersServer {
	return &UsersServer{log: logger.Named("UsersServer"), ctrl: ctrl}
}

func (s *UsersServer) FetchDefaults(ctx context.Context, req *connect.Request[cc.Empty]) (
	*connect.Response[cc.Defaults], error) {
	log := s.log.Named("FetchDefaults")
	log.Debug("Request received", zap.Any("req", req.Msg))

	defaults, err := core.Config()
	if err != nil {
		s.log.Error("Failed get config", zap.Error(err))
		return nil, fmt.Errorf("failed to fetch defaults: %w", err)
	}

	resp := connect.NewResponse[cc.Defaults](defaults)

	return resp, nil
}

func (s *UsersServer) Resolve(ctx context.Context, req *connect.Request[cc.Users]) (*connect.Response[cc.Users], error) {
	log := s.log.Named("Resolve")
	log.Debug("Request received", zap.Any("req", req.Msg))

	requestor := ctx.Value(core.ChatAccount).(string)

	conf, err := core.Config()
	if err != nil {
		s.log.Error("Failed get config", zap.Error(err))
		return nil, fmt.Errorf("failed to fetch defaults: %w", err)
	}

	uuids := append(conf.Admins, requestor)
	for _, user := range req.Msg.GetUsers() {
		uuids = append(uuids, user.GetUuid())
	}

	uuids = core.Unique(uuids)

	res, err := s.ctrl.Resolve(ctx, uuids)
	if err != nil {
		return nil, err
	}

	resp := connect.NewResponse[cc.Users](&cc.Users{
		Users: res,
	})
	return resp, nil
}

func (s *UsersServer) Me(ctx context.Context, req *connect.Request[cc.Empty]) (*connect.Response[cc.User], error) {
	log := s.log.Named("Me")
	log.Debug("Request received", zap.Any("req", req.Msg))

	requestor := ctx.Value(core.ChatAccount).(string)

	res, err := s.ctrl.Resolve(ctx, []string{requestor})
	if err != nil {
		return nil, err
	}

	for _, user := range res {
		if user.Uuid == requestor {
			resp := connect.NewResponse[cc.User](user)
			return resp, nil
		}
	}

	return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("undexpected: User not found"))
}

func (s *UsersServer) GetMembers(ctx context.Context, req *connect.Request[cc.Empty]) (*connect.Response[cc.Users], error) {
	log := s.log.Named("Me")
	log.Debug("Request received", zap.Any("req", req.Msg))

	//requestor := ctx.Value(core.ChatAccount).(string)

	members, err := s.ctrl.GetMembers(ctx)
	if err != nil {
		return nil, err
	}

	resp := connect.NewResponse[cc.Users](&cc.Users{
		Users: members,
	})

	return resp, nil
}
