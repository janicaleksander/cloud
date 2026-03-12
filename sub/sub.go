package sub

import (
	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/event"
	"github.com/janicaleksander/cloud/utils"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Subscriber struct {
	ID      uuid.UUID
	Channel *amqp.Channel
	Queue   amqp.Queue
	Type    any
}

func NewSubscriber[T event.Event](conn *amqp.Connection) (*Subscriber, error) {
	var sample T
	name := utils.GetTypeName(sample)

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		name, false, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	return &Subscriber{
		ID:      uuid.New(),
		Channel: ch,
		Queue:   q,
		Type:    name,
	}, nil
}

func (s *Subscriber) Consume() (<-chan amqp.Delivery, error) {
	msgs, err := s.Channel.Consume(
		s.Queue.Name,
		"",
		true, // setting auto ack to automatically send ack after received message
		false,
		false,
		false,

		nil)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}
