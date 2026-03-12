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

// Drugi konsument zdarzenia typu 1
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

	// Type1Event subscriber
	s21, err := sub.NewSubscriber[event.Type1Event](conn)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	msgs21, err := s21.Consume()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	go func() {
		for msg := range msgs21 {
			slog.Info(fmt.Sprintf(
				"Received a message: %s, from queue: %s, msg type %s, by consumer %s",
				msg.Body,
				s21.Queue.Name,
				s21.Type,
				s21.ID))
		}
	}()

	<-done
}
