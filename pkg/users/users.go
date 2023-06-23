package users

import (
	"context"

	"github.com/bufbuild/connect-go"
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

	defaults := &cc.Defaults{
		Gateways: []string{"whmcs", "email", "telegram"},
		Admins: []*cc.User{
			{
				Uuid: "0", Title: "NoCloud",
			},
		},
	}
	resp := connect.NewResponse[cc.Defaults](defaults)

	return resp, nil
}

func (s *UsersServer) Resolve(ctx context.Context, req *connect.Request[cc.Users]) (*connect.Response[cc.Users], error) {
	log := s.log.Named("Resolve")
	log.Debug("Request received", zap.Any("req", req.Msg))

	requestor := ctx.Value(core.ChatAccount).(string)

	res, err := s.ctrl.Resolve(ctx, []string{requestor})
	if err != nil {
		return nil, err
	}

	resp := connect.NewResponse[cc.Users](&cc.Users{
		Users: res,
	})
	return resp, nil
}
