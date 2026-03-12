package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/janicaleksander/cloud/event"
	"github.com/janicaleksander/cloud/pub"
	"github.com/janicaleksander/cloud/sub"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

// 1 konsumenta zdarzenia typu 3, który po przetworzeniu zdarzenia wygeneruje zdarzenie typu 4 i je opublikuje
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

	// Type3Event subscriber - generates and publishes Type4Event
	s13, err := sub.NewSubscriber[event.Type3Event](conn)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	msgs13, err := s13.Consume()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	p14, err := pub.NewPublisher(conn)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	go func() {
		for msg := range msgs13 {
			slog.Info(fmt.Sprintf(
				"Received a message: %s, from queue: %s, msg type %s, by consumer %s",
				msg.Body,
				s13.Queue.Name,
				s13.Type,
				s13.ID))
			// Generate and publish Type4Event
			p14.Publish(event.NewType4Event())
		}
	}()

	<-done
}
