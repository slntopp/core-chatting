package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"github.com/slntopp/core-chatting/cc"
	"go.uber.org/zap"
)

type PubSub struct {
	log *zap.Logger

	ch *amqp091.Channel
}

func NewPubSub(logger *zap.Logger, conn *amqp091.Connection) *PubSub {
	log := logger.Named("PubSub")

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to create channel", zap.Error(err))
	}

	return &PubSub{log: log, ch: ch}
}

func (s *PubSub) Pub(ctx context.Context, id string, event *cc.Event) {
	log := s.log.Named(fmt.Sprintf("Pub-%s", id))

	q, err := s.ch.QueueDeclare(
		id,
		false,
		true,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Error("Failed to get queue", zap.Error(err))
		return
	}

	marshal, err := json.Marshal(event)

	if err != nil {
		log.Error("Failed to marshal event", zap.Error(err))
		return
	}

	err = s.ch.PublishWithContext(
		ctx,
		"",
		q.Name,
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

func (s *PubSub) Sub(id string) <-chan amqp091.Delivery {
	log := s.log.Named(fmt.Sprintf("Sub-%s", id))

	q, err := s.ch.QueueDeclare(
		id,
		false,
		true,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Error("Failed to get queue", zap.Error(err))
		return nil
	}

	msgs, err := s.ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Error("Failed to get consume chan", zap.Error(err))
		return nil
	}

	return msgs
}
