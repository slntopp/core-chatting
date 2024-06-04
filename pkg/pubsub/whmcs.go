package pubsub

import (
	"context"

	"github.com/rabbitmq/amqp091-go"
	"github.com/slntopp/core-chatting/cc"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (s *PubSub) PubWhmcs(ctx context.Context, event *cc.Event) {
	log := s.log.Named("Pub-Whmcs")
	queueName := "whmcs_tickets"

	channel, err := s.conn.Channel()

	if err != nil {
		log.Error("Failed to init channel", zap.Error(err))
		return
	}
	defer channel.Close()

	queue, err := channel.QueueDeclare(queueName, true, false, false, false, nil)
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
