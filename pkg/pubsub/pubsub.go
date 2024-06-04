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

	conn *amqp091.Connection
}

func NewPubSub(logger *zap.Logger, conn *amqp091.Connection) *PubSub {
	log := logger.Named("PubSub")
	return &PubSub{log: log, conn: conn}
}

func (s *PubSub) Pub(ctx context.Context, id string, event *cc.Event) {
	log := s.log.Named(fmt.Sprintf("Pub-%s", id))

	exchange := fmt.Sprintf("cc.user.%s", id)

	channel, err := s.conn.Channel()

	if err != nil {
		log.Error("Failed to init channel", zap.Error(err))
		return
	}
	defer channel.Close()

	err = channel.ExchangeDeclare(
		exchange,
		"fanout",
		true,
		false,
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

	err = channel.PublishWithContext(
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

	channel, err := s.conn.Channel()

	if err != nil {
		log.Error("Failed to init channel", zap.Error(err))
		return nil, nil, err
	}

	err = channel.ExchangeDeclare(
		exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Error("failed to create exchange", zap.Error(err))
		return nil, nil, err
	}

	q, err := channel.QueueDeclare(
		"",
		true,
		false,
		true,
		false,
		nil,
	)

	if err != nil {
		log.Error("Failed to get queue", zap.Error(err))
		return nil, nil, err
	}

	queueTerminator := func() error {
		_, err := channel.QueueDelete(q.Name, false, false, false)
		channel.Close()
		return err
	}

	err = channel.QueueBind(
		q.Name,
		"",
		exchange,
		false,
		nil,
	)

	if err != nil {
		log.Error("bind queue", zap.Error(err))
		return nil, nil, err
	}

	msgs, err := channel.Consume(
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
