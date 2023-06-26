package stream

import (
	"context"
	"encoding/json"
	"github.com/bufbuild/connect-go"
	"github.com/slntopp/core-chatting/cc"
	"github.com/slntopp/core-chatting/pkg/core"
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
	requestor := ctx.Value(core.ChatAccount).(string)

	log := s.log.Named("Stream").Named(requestor)

	res, err := s.ctrl.Resolve(ctx, []string{requestor})
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	for _, user := range res {
		if user.Uuid == requestor {
			goto start_stream
		}
	}

	log.Error("Failed to resolve user")

	return nil

start_stream:

	msgs := s.ps.Sub(requestor)

	var event = &cc.Event{}

	for msg := range msgs {
		err := json.Unmarshal(msg.Body, event)
		if err != nil {
			return err
		}

		err = serverStream.Send(event)
		if err != nil {
			return err
		}
	}

	return nil
}
