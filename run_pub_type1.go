package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/janicaleksander/cloud/event"
	"github.com/janicaleksander/cloud/pub"
	"github.com/janicaleksander/cloud/utils"
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
	exchange := os.Getenv("EXCHANGE_NAME")

	// Create 1 Type1 publisher
	publisher, err := pub.NewPublisher(conn, exchange)
	if err != nil {
		slog.Error("Spawning publisher error ", err.Error())
		return
	}
	defer publisher.Channel.Close()

	// Start Type1 publisher
	go func() {
		for {
			ctx, cancel := context.WithTimeout(context.Background(), publishTimeout)
			if err := publisher.Publish(ctx, event.NewType1Event()); err != nil {
				slog.Error("Failed to publish: " + err.Error())
			}
			cancel()
			time.Sleep(utils.Delay(type1Delay)())
		}
	}()

	slog.Info("Type1 publisher started (2s delay)")

	done := make(chan struct{})
	<-done
}
