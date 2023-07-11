package pubsub

import (
	"context"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"github.com/slntopp/core-chatting/cc"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (s *PubSub) PubGateway(ctx context.Context, id string, event *cc.Event, gateways []string) {
	log := s.log.Named(fmt.Sprintf("Pub-Gateway.%s", id))
	for _, gateway := range gateways {
		exchangeName := fmt.Sprintf("cc.gateway.%s", gateway)

		err := s.ch.ExchangeDeclare(
			exchangeName,
			"topic",
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

		queueName := fmt.Sprintf("cc.gateway.%s.%s", gateway, id)

		queue, err := s.ch.QueueDeclare(queueName, true, false, true, false, nil)
		if err != nil {
			log.Error("failed to create exchange", zap.Error(err))
			return
		}

		err = s.ch.QueueBind(queue.Name, id, exchangeName, false, nil)
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
			exchangeName,
			id,
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
}
