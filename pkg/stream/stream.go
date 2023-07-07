package stream

import (
	"context"
	"errors"
	"github.com/bufbuild/connect-go"
	"net/http"
	"time"

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
	userCtrl *graph.UsersController
	msgCtrl  *graph.MessagesController
	ps       *pubsub.PubSub
	key      []byte
}

func NewStreamServer(logger *zap.Logger, userCtrl *graph.UsersController, msgCtrl *graph.MessagesController, ps *pubsub.PubSub, key []byte) *StreamServer {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		Subprotocols: []string{"Bearer"},
	}

	return &StreamServer{
		log:      logger.Named("StreamService"),
		userCtrl: userCtrl,
		msgCtrl:  msgCtrl,
		ps:       ps,
		upgrader: upgrader,
		key:      key,
	}
}

func (s *StreamServer) Stream(ctx context.Context, req *connect.Request[cc.Empty], serverStream *connect.ServerStream[cc.Event]) error {
	requestor := ctx.Value(core.ChatAccount).(string)

	log := s.log.Named("Stream").Named(requestor)

	defer log.Debug("Stream closed")

	res, err := s.userCtrl.Resolve(ctx, []string{requestor})
	if err != nil {
		log.Error("Failed to resolve user", zap.Error(err))
		return err
	}

	for _, user := range res {
		if user.Uuid == requestor {
			goto start_stream
		}
	}

	log.Error("Failed to resolve user")

	return connect.NewError(connect.CodeUnauthenticated, errors.New("failed to resolve user"))
start_stream:

	log.Info("Start stream", zap.String("user", requestor))

	msgs, queueTerminator, err := s.ps.Sub(requestor)

	defer func() {
		queueError := queueTerminator()
		if queueError != nil {
			log.Error("Failed to delete queue", zap.Error(queueError))
		}
	}()

	if err != nil {
		return connect.NewError(connect.CodeInternal, err)
	}

	var event = &cc.Event{}

	ticker := time.NewTicker(10 * time.Second)

	for {
		select {
		case msg := <-msgs:
			err := proto.Unmarshal(msg.Body, event)
			if err != nil {
				log.Error("Failed to unmarshal event", zap.Error(err))
				continue
			}

			log.Info("Receive message", zap.Any("event", event))

			if event.GetType() == cc.EventType_MESSAGE_SENT {
				message := event.GetMsg()
				newMessage, err := s.msgCtrl.Read(ctx, message, requestor)
				if err != nil {
					log.Error("Failed to update readers", zap.Error(err))
				}
				event.Item = &cc.Event_Msg{Msg: newMessage}
			}

			err = serverStream.Send(event)
			if err != nil {
				log.Error("Error writing message", zap.Error(err))
				return err
			}

			log.Debug("Processed event successfully")
		case <-ticker.C:
			err := serverStream.Send(&cc.Event{
				Type: cc.EventType_PING,
			})
			if err != nil {
				return nil
			}
		}
	}
}
