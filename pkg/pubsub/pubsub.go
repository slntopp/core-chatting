package pubsub

import (
	"context"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
	"github.com/slntopp/core-chatting/cc"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
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

	exchange := fmt.Sprintf("cc.user.%s", id)

	err := s.ch.ExchangeDeclare(
		exchange,
		"fanout",
		false,
		true,
		false,
		false,
		nil,
	)

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
		exchange,
		"",
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

func (s *PubSub) Sub(id string) (<-chan amqp091.Delivery, func() error, error) {
	log := s.log.Named(fmt.Sprintf("Sub-%s", id))

	exchange := fmt.Sprintf("cc.user.%s", id)

	err := s.ch.ExchangeDeclare(
		exchange,
		"fanout",
		false,
		true,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Error("failed to create exchange", zap.Error(err))
		return nil, nil, err
	}

	q, err := s.ch.QueueDeclare(
		"",
		false,
		true,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Error("Failed to get queue", zap.Error(err))
		return nil, nil, err
	}

	queueTerminator := func() error {
		_, err := s.ch.QueueDelete(q.Name, false, false, false)
		return err
	}

	err = s.ch.QueueBind(
		q.Name,
		"",
		exchange,
		false,
		nil,
	)

	msgs, err := s.ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Error("Failed to get consume chan", zap.Error(err))
		return nil, nil, err
	}

	return msgs, queueTerminator, nil
}
