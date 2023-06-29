package stream

import (
	"context"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/slntopp/core-chatting/cc"
	"github.com/slntopp/core-chatting/pkg/core"
	"github.com/slntopp/core-chatting/pkg/graph"
	"github.com/slntopp/core-chatting/pkg/pubsub"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type StreamServer struct {
	log *zap.Logger

	upgrader websocket.Upgrader
	ctrl     *graph.UsersController
	ps       *pubsub.PubSub
}

func NewStreamServer(logger *zap.Logger, ctrl *graph.UsersController, ps *pubsub.PubSub) *StreamServer {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	return &StreamServer{
		log:      logger.Named("StreamService"),
		ctrl:     ctrl,
		ps:       ps,
		upgrader: upgrader,
	}
}

func (s *StreamServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestor := r.Header.Get(string(core.ChatAccount))

	connection, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	defer connection.Close()

	ctx := context.Background()

	log := s.log.Named("Stream").Named(requestor)

	res, err := s.ctrl.Resolve(ctx, []string{requestor})
	if err != nil {
		return
	}

	for _, user := range res {
		if user.Uuid == requestor {
			goto start_stream
		}
	}

	log.Error("Failed to resolve user")

	return
start_stream:

	log.Info("Start stream", zap.String("user", requestor))

	connection.WriteMessage(websocket.TextMessage, []byte("Establised connection"))

	msgs := s.ps.Sub(requestor)

	var event = &cc.Event{}

	for msg := range msgs {
		err := proto.Unmarshal(msg.Body, event)
		if err != nil {
			log.Error("Failed to unmarshal event", zap.Error(err))
		}

		log.Info("Receive message", zap.Any("event", event))

		err = connection.WriteMessage(websocket.BinaryMessage, msg.Body)
		if err != nil {
			log.Error("Error", zap.Error(err))
		}

		err = msg.Ack(false)
		if err != nil {
			log.Error("Failed to ack msg", zap.Error(err))
		}
	}
}
