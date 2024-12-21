package chats

import (
	"context"
	"github.com/rabbitmq/amqp091-go"
	"github.com/slntopp/core-chatting/cc"
	"github.com/slntopp/core-chatting/pkg/pubsub"
	"github.com/slntopp/nocloud-proto/events"
	"github.com/slntopp/nocloud/pkg/nocloud/schema"
	sc "github.com/slntopp/nocloud/pkg/settings/client"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
	"sync"
	"time"
)

func (s *ChatsServer) sendCloseChatMessage(ctx context.Context, log *zap.Logger, chat *cc.Chat, msg *cc.Message) error {
	message, err := s.msgsCtrl.Send(ctx, msg)
	if err != nil {
		return err
	}

	go pubsub.HandleNotifyMessage(ctx, log, s.ps, message, chat, cc.EventType_MESSAGE_SENT)

	if s.whmcsTickets && chat.GetDepartment() != "openai" {
		go s.ps.PubWhmcs(ctx, &cc.Event{
			Type: cc.EventType_MESSAGE_SENT,
			Item: &cc.Event_Msg{Msg: message},
		})
	}

	return nil
}

func (s *ChatsServer) CloseInactiveChats(ctx context.Context, log *zap.Logger, conf TicketsSettingsConf, eventPublisher func(ctx context.Context, event *events.Event)) error {
	log = log.Named("CloseInactiveChats")

	chats, err := s.ctrl.List(ctx, schema.ROOT_ACCOUNT_KEY)
	if err != nil {
		return err
	}

	now := time.Now().Unix()
	closeAfterSeconds := int64(conf.CloseInactiveChatsAfterHours * 60 * 60)
	for _, chat := range chats {
		if chat.Status != cc.Status_ANSWERED {
			continue
		}
		if chat.Meta == nil || chat.GetMeta().LastMessage == nil {
			continue
		}
		last := chat.GetMeta().GetLastMessage()
		if last.Sent == 0 {
			continue
		}
		if now-(last.Sent/1000) < closeAfterSeconds {
			continue
		}
		chat.Status = cc.Status_CLOSE
		if _, err = s.ctrl.Update(ctx, chat); err != nil {
			return err
		}
		go eventPublisher(ctx, &events.Event{
			Key:  "inactive_chat_closed",
			Type: "email",
			Data: map[string]*structpb.Value{
				"subject": structpb.NewStringValue(chat.GetTopic()),
			},
			Uuid: chat.GetOwner(),
		})
		if err = s.sendCloseChatMessage(ctx, log, chat, &cc.Message{
			Sent:    time.Now().UnixMilli(),
			Sender:  schema.ROOT_ACCOUNT_KEY,
			Content: conf.CloseMessageContent,
			Chat:    &chat.Uuid,
		}); err != nil {
			return err
		}
	}

	return nil
}

func setupEventPublisher(_ context.Context, log *zap.Logger, rbmq *amqp091.Connection) func(ctx context.Context, event *events.Event) {
	ch, err := rbmq.Channel()
	if err != nil {
		log.Fatal("Failed to open a channel", zap.Error(err))
	}
	return func(ctx context.Context, event *events.Event) {
		queue, _ := ch.QueueDeclare(
			"events",
			true, false, false, true, nil,
		)
		body, err := proto.Marshal(event)
		if err != nil {
			log.Error("Error while marshalling record", zap.Error(err))
			return
		}
		if err = ch.PublishWithContext(ctx, "", queue.Name, false, false, amqp091.Publishing{
			ContentType: "text/plain", Body: body,
		}); err != nil {
			log.Error("Error while publishing event", zap.Error(err))
		}
	}
}

func (s *ChatsServer) CloseInactiveChatsRoutine(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	log := s.log.Named("CloseInactiveChatsRoutine")
	eventPublisher := setupEventPublisher(ctx, log, s.conn)

start:

	ticketsConf := MakeTicketsConf(ctx, log, &s.settingsClient)

	upd := make(chan bool, 1)
	go sc.Subscribe([]string{ticketsKey}, upd)

	log.Info("Got Configuration", zap.Any("conf", ticketsConf))

	ticker := time.NewTicker(time.Second * time.Duration(ticketsConf.RoutineFrequency))
	tick := time.Now()
	for {
		log.Info("Entering new Iteration", zap.Time("ts", tick))
		if err := s.CloseInactiveChats(ctx, log, ticketsConf, eventPublisher); err != nil {
			log.Error("Error while closing inactive chats", zap.Error(err))
		}
		select {
		case <-ctx.Done():
			log.Info("Context is done. Quitting")
			return
		case tick = <-ticker.C:
			continue
		case <-upd:
			log.Info("New Configuration Received, restarting Routine")
			goto start
		}
	}
}
