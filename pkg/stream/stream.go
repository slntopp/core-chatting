package stream

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/slntopp/core-chatting/cc"
	"github.com/slntopp/core-chatting/pkg/core"

	"connectrpc.com/connect"
	"github.com/slntopp/core-chatting/pkg/graph"
	"github.com/slntopp/core-chatting/pkg/pubsub"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type StreamServer struct {
	log *zap.Logger

	upgrader  websocket.Upgrader
	userCtrl  *graph.UsersController
	msgCtrl   *graph.MessagesController
	chatsCtrl *graph.ChatsController
	ps        *pubsub.PubSub
	key       []byte
}

func NewStreamServer(logger *zap.Logger, userCtrl *graph.UsersController, msgCtrl *graph.MessagesController, chatsCtrl *graph.ChatsController, ps *pubsub.PubSub, key []byte) *StreamServer {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		Subprotocols: []string{"Bearer"},
	}

	return &StreamServer{
		log:       logger.Named("StreamService"),
		userCtrl:  userCtrl,
		msgCtrl:   msgCtrl,
		chatsCtrl: chatsCtrl,
		ps:        ps,
		upgrader:  upgrader,
		key:       key,
	}
}

func (s *StreamServer) Stream(ctx context.Context, req *connect.Request[cc.StreamRequest], serverStream *connect.ServerStream[cc.Event]) error {
	requestor := ctx.Value(core.ChatAccount).(string)

	log := s.log.Named("Stream").Named(requestor)

	defer log.Debug("Stream closed")

	res, err := s.userCtrl.Resolve(ctx, []string{requestor})
	if err != nil {
		log.Error("Failed to resolve user", zap.Error(err))
		return err
	}

	var streamUser *cc.User

	for _, user := range res {
		if user.Uuid == requestor {
			streamUser = user
			goto start_stream
		}
	}

	log.Error("Failed to resolve user")

	return connect.NewError(connect.CodeUnauthenticated, errors.New("failed to resolve user"))
start_stream:

	if streamUser.GetCcIsBot() {
		log.Debug("Stream requestor is a bot", zap.String("username", streamUser.GetCcUsername()), zap.Any("commands", req.Msg.GetCommands()))
		err := s.userCtrl.UpdateCommands(ctx, streamUser.GetUuid(), req.Msg.GetCommands())
		if err != nil {
			log.Error("Failed to update commands in bot", zap.Error(err))
			return err
		}
	}

	log.Info("Start stream", zap.String("user", requestor))

	msgs, queueTerminator, err := s.ps.Sub(requestor)

	if err != nil {
		return connect.NewError(connect.CodeInternal, err)
	}

	defer func() {
		log.Debug("Queue termination")
		queueError := queueTerminator()
		if queueError != nil {
			log.Error("Failed to delete queue", zap.Error(queueError))
		}
	}()

	var event = &cc.Event{}

	ticker := time.NewTicker(10 * time.Second)

	for {
		select {
		case msg, opened := <-msgs:
			if !opened {
				log.Warn("Channel was closed")
				return nil
			}

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
					return err
				}
				chat, err := s.chatsCtrl.Get(ctx, message.GetChat(), requestor)
				if err != nil {
					log.Error("Failed to get chat", zap.Error(err))
					return err
				}

				event.Item = &cc.Event_Msg{Msg: newMessage}

				readEvent := &cc.Event{
					Type: cc.EventType_CHAT_READ,
					Item: &cc.Event_Chat{
						Chat: chat,
					},
				}
				err = serverStream.Send(readEvent)
				if err != nil {
					log.Error("Failed to send read event", zap.Error(err))
					return err
				}
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
				log.Error("Failed to send message", zap.Error(err))
				return nil
			}
		}
	}
}
