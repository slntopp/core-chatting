package stream

import (
	"context"
	"github.com/bufbuild/connect-go"
	"github.com/slntopp/core-chatting/cc"
	"github.com/slntopp/core-chatting/pkg/graph"
	"github.com/slntopp/core-chatting/pkg/pubsub"
	"go.uber.org/zap"
)

type StreamServer struct {
	log *zap.Logger

	ctrl *graph.UsersController
	ps   *pubsub.PubSub
}

func NewStreamServer(logger *zap.Logger, ctrl *graph.UsersController, ps *pubsub.PubSub) *StreamServer {
	return &StreamServer{
		log:  logger.Named("StreamService"),
		ctrl: ctrl,
		ps:   ps,
	}
}

func (s *StreamServer) Stream(ctx context.Context, req *connect.Request[cc.Empty], serverStream *connect.ServerStream[cc.Event]) error {
	return nil
}
