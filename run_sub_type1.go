package main

import (
	"log/slog"
	"os"

	"github.com/janicaleksander/cloud/event"
	"github.com/janicaleksander/cloud/messaging"
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

	// Create 1 Type1Event subscriber
	messaging.Subscribe[event.Type1Event](conn, messaging.ProcessMessages)

	slog.Info("Type1 subscriber started")
	done := make(chan struct{})
	<-done
}
