package sub

import (
	"log/slog"

	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/event"
	"github.com/janicaleksander/cloud/utils"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Subscriber struct {
	ID           uuid.UUID
	ExchangeName string
	Channel      *amqp.Channel
	Queue        amqp.Queue // queue name is the same as type sending through this queue (routing key == queue name)
}

func NewSubscriber[T event.Event](conn *amqp.Connection, exchangeName string) (*Subscriber, error) {
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
	s := &Subscriber{
		ID:           uuid.New(),
		ExchangeName: exchangeName,
		Channel:      ch,
		Queue:        q,
	}
	err = s.Bind()
	if err != nil {
		return nil, err
	}
	defer func() {
		slog.Info("Creating subscriber with ID: " + s.ID.String() + ", for queue: " + s.Queue.Name)
	}()
	return s, nil
}
func (s *Subscriber) Bind() error {
	return s.Channel.QueueBind(
		s.Queue.Name,   // queue name
		s.Queue.Name,   // routing key
		s.ExchangeName, //exchange name
		false,
		nil)

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
