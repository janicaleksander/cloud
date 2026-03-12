package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/janicaleksander/cloud/event"
	"github.com/janicaleksander/cloud/messaging"
	"github.com/janicaleksander/cloud/pub"
	"github.com/janicaleksander/cloud/sub"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

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
	messaging.Subscribe[event.Type1Event](conn, messaging.ProcessMessages)
	messaging.Subscribe[event.Type1Event](conn, messaging.ProcessMessages)
	messaging.Subscribe[event.Type2Event](conn, messaging.ProcessMessages)
	messaging.Subscribe[event.Type4Event](conn, messaging.ProcessMessages)

	// Type3Event -> publikuje Type4Event
	p14, err := pub.NewPublisher(conn, os.Getenv("EXCHANGE_NAME"))
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer p14.Channel.Close()

	messaging.Subscribe[event.Type3Event](conn, func(msgs <-chan amqp.Delivery, s *sub.Subscriber) {
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

	done := make(chan struct{})
	<-done
}
