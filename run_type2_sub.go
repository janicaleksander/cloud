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

// 1 konsumenta zdarzenia typu 2
func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Can't load .env file")
		return
	}
	conn, err := amqp.Dial(os.Getenv("AMQP_URL"))
	if err != nil {
		slog.Error("Can't connect to amqp")
		return
	}
	defer conn.Close()
	done := make(chan struct{})

	// Type2Event subscriber
	s12, err := sub.NewSubscriber[event.Type2Event](conn)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	msgs12, err := s12.Consume()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	go func() {
		for msg := range msgs12 {
			slog.Info(fmt.Sprintf(
				"Received a message: %s, from queue: %s, msg type %s, by consumer %s",
				msg.Body,
				s12.Queue.Name,
				s12.Type,
				s12.ID))
		}
	}()

	<-done
}
