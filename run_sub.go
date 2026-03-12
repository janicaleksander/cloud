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
			subscriber.Type,
			subscriber.ID))
	}
}

func mustSubscribe[T event.Event](conn *amqp.Connection, handler func(<-chan amqp.Delivery, *sub.Subscriber)) {
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

	// 2x Type1Event, 1x Type2Event, 1x Type4Event
	mustSubscribe[event.Type1Event](conn, processMessages)
	mustSubscribe[event.Type1Event](conn, processMessages)
	mustSubscribe[event.Type2Event](conn, processMessages)
	mustSubscribe[event.Type4Event](conn, processMessages)

	// Type3Event -> publikuje Type4Event
	p14, err := pub.NewPublisher(conn, os.Getenv("EXCHANGE_NAME"))
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer p14.Channel.Close()

	mustSubscribe[event.Type3Event](conn, func(msgs <-chan amqp.Delivery, s *sub.Subscriber) {
		for msg := range msgs {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			slog.Info(
				fmt.Sprintf("Received a message: %s, from queue: %s, msg type %s, by consumer %s",
					msg.Body,
					s.Queue.Name,
					s.Type,
					s.ID,
				),
			)
			if err := p14.Publish(ctx, event.NewType4Event()); err != nil {
				slog.Error("Failed to publish Type4Event: " + err.Error())
			}
			cancel()
		}
	})

	done := make(chan struct{})
	<-done
}
