package pubsub

import (
	"context"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"github.com/slntopp/core-chatting/cc"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (s *PubSub) PubTicket(ctx context.Context, event *cc.Event) {
	log := s.log.Named(fmt.Sprintf("PubTicket"))

	queue, err := s.ch.QueueDeclare("tickets", true, false, false, false, nil)
	if err != nil {
		log.Error("failed to create exchange", zap.Error(err))
		return
	}

	marshal, err := proto.Marshal(event)
	if err != nil {
		log.Error("Failed to marshal event", zap.Error(err))
		return
	}

	log.Debug("Publish event", zap.Any("event", event))

	err = s.ch.PublishWithContext(
		ctx,
		"",
		queue.Name,
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        marshal,
		},
	)
	if err != nil {
		log.Error("Failed to publish event", zap.Error(err))
		return
	}
}
