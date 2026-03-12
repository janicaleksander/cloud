package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/janicaleksander/cloud/event"
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

	// Create 1 Type4Event subscriber
	subscribe[event.Type4Event](conn, processMessages)

	slog.Info("Type4 subscriber started")

	done := make(chan struct{})
	<-done
}
