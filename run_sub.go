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

/*
2 konsumentów zdarzenia typu 1
1 konsumenta zdarzenia typu 2
1 konsumenta zdarzenia typu 3, który po przetworzeniu zdarzenia wygeneruje zdarzenie typu 4 i je opublikuje.
1 konsumenta zdarzenia typu 4.
*/

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

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Can't load .env file")
		return
	}
	conn, err := amqp.Dial(os.Getenv("AMQP_URL"))
	if err != nil {
		slog.Error("Can't connect to amqp ")
		return
	}
	defer conn.Close()
	done := make(chan struct{})

	// First Type1Event subscriber
	s11, err := sub.NewSubscriber[event.Type1Event](conn)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	msgs11, err := s11.Consume()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	go func() {
		processMessages(msgs11, s11)
	}()

	// Second Type1Event subscriber
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
		processMessages(msgs21, s21)
	}()

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
		processMessages(msgs12, s12)
	}()

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
			slog.Info(fmt.Sprintf("Received a message: %s, from queue: %s, msg type %s, by consumer %s", msg.Body, s13.Queue.Name, s13.Type, s13.ID))
			// Generate and publish Type4Event
			p14.Publish(event.NewType4Event())
		}
	}()

	// Type4Event subscriber
	s14, err := sub.NewSubscriber[event.Type4Event](conn)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	msgs14, err := s14.Consume()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	go func() {
		processMessages(msgs14, s14)
	}()

	<-done
}
