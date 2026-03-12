package messaging

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/janicaleksander/cloud/event"
	"github.com/janicaleksander/cloud/sub"
	amqp "github.com/rabbitmq/amqp091-go"
)

func Subscribe[T event.Event](conn *amqp.Connection, handler func(<-chan amqp.Delivery, *sub.Subscriber)) {
	s, err := sub.NewSubscriber[T](conn, os.Getenv("EXCHANGE_NAME"))
	if err != nil {
		slog.Error(err.Error())
		return
	}
	msgs, err := s.Consume()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	go handler(msgs, s)
}

func ProcessMessages(delivery <-chan amqp.Delivery, subscriber *sub.Subscriber) {
	for msg := range delivery {
		slog.Info(fmt.Sprintf(
			"Received a message: %s, from queue: %s, msg type %s, by consumer %s",
			msg.Body,
			subscriber.Queue.Name,
			subscriber.Queue.Name,
			subscriber.ID))
	}
}
