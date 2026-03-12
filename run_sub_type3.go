package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/janicaleksander/cloud/event"
	"github.com/janicaleksander/cloud/pub"
	"github.com/janicaleksander/cloud/sub"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func processMessages(delivery <-chan amqp.Delivery, subscriber *sub.Subscriber) {
	for msg := range delivery {
		slog.Info(fmt.Sprintf(
			"Received a message: %s, from queue: %s, msg type %s, by consumer %s",
			msg.Body,
			subscriber.Queue.Name,
			subscriber.Queue.Name,
			subscriber.ID))
	}
}

func subscribe[T event.Event](conn *amqp.Connection, handler func(<-chan amqp.Delivery, *sub.Subscriber)) {
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
func main() {
	if err := godotenv.Load(); err != nil {
		slog.Error("Can't load .env file")
		return
	}
	conn, err := amqp.Dial(os.Getenv("AMQP_URL"))
	if err != nil {
		slog.Error("Can't connect to amqp")
		return
	}
	defer conn.Close()

	// Create publisher for Type4Event
	p14, err := pub.NewPublisher(conn, os.Getenv("EXCHANGE_NAME"))
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer p14.Channel.Close()

	// Subscribe to Type3Event and publish Type4Event
	subscribe[event.Type3Event](conn, func(msgs <-chan amqp.Delivery, s *sub.Subscriber) {
		for msg := range msgs {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			slog.Info(
				fmt.Sprintf("Received a message: %s, from queue: %s, msg type %s, by consumer %s",
					msg.Body,
					s.Queue.Name,
					s.Queue.Name,
					s.ID,
				),
			)
			if err := p14.Publish(ctx, event.NewType4Event()); err != nil {
				slog.Error("Failed to publish Type4Event: " + err.Error())
			}
			cancel()
		}
	})

	slog.Info("Type3 subscriber started (processes Type3Event -> publishes Type4Event)")

	done := make(chan struct{})
	<-done
}
