package stream

import (
	"context"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/slntopp/core-chatting/cc"
	"github.com/slntopp/core-chatting/pkg/core"

	"github.com/slntopp/core-chatting/pkg/core/auth"
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
	key      []byte
}

func NewStreamServer(logger *zap.Logger, ctrl *graph.UsersController, ps *pubsub.PubSub, key []byte) *StreamServer {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		Subprotocols: []string{"Bearer"},
	}

	return &StreamServer{
		log:      logger.Named("StreamService"),
		ctrl:     ctrl,
		ps:       ps,
		upgrader: upgrader,
		key:      key,
	}
}

func (s *StreamServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log := s.log.Named("Stream")
	log.Debug("New connection", zap.Any("headers", r.Header))

	connection, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Warn("Failed to upgrade connection", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer func() {
		log.Debug("Closing connection")
		log.Debug("Connection closed", zap.Error(connection.Close()))
	}()

	log.Debug("Negotiated Sub Protocol", zap.String("protocol", connection.Subprotocol()))

	protocol := r.Header.Get("Sec-Websocket-Protocol")
	protocols := strings.Split(protocol, ", ")

	if len(protocols) != 2 {
		log.Warn("Invalid protocol")
		connection.WriteMessage(websocket.CloseMessage, []byte("Invalid protocol"))
		return
	}

	claims, err := auth.ValidateToken(s.key, protocols[1])
	if err != nil {
		log.Debug("Failed to validate token", zap.Error(err))
		connection.WriteMessage(websocket.CloseMessage, []byte("Failed to validate token"))
		return
	}

	requestor := claims[core.JWT_ACCOUNT_CLAIM].(string)

	ctx := context.Background()

	log = log.Named(requestor)

	res, err := s.ctrl.Resolve(ctx, []string{requestor})
	if err != nil {
		log.Error("Failed to resolve user", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
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

	msgs, err := s.ps.Sub(requestor)

	if err != nil {
		connection.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		return
	}

	var event = &cc.Event{}

	for msg := range msgs {
		err := proto.Unmarshal(msg.Body, event)
		if err != nil {
			log.Error("Failed to unmarshal event", zap.Error(err))
		}

		log.Info("Receive message", zap.Any("event", event))

		err = connection.WriteMessage(websocket.BinaryMessage, msg.Body)
		if err != nil {
			log.Error("Error writing message", zap.Error(err))
			return
		}

		err = msg.Ack(false)
		if err != nil {
			log.Error("Failed to ack msg", zap.Error(err))
		}

		log.Debug("Processed event successfully")
	}
}
