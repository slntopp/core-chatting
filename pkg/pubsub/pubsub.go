package pubsub

import "github.com/rabbitmq/amqp091-go"

type PubSub struct {
	conn *amqp091.Connection
}

func (s *PubSub) Pub(id string) {

}

func (s *PubSub) Sub(id string) {

}
