package main

import (
	"log/slog"
	"math/rand"
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

	// Create Type3 publisher
	publisher, err := pub.NewPublisher(conn, exchange)
	if err != nil {
		slog.Error("Spawning publisher error ", err.Error())
		return
	}
	defer publisher.Channel.Close()

	// Start Type3 publisher with random delay (1-5 seconds)
	startPublisher(publisher, func() any { return event.NewType3Event() }, utils.Delay(time.Second*time.Duration(rand.Intn(randomDelayRange)+1)))

	slog.Info("Type3 publisher started (random 1-5s delay)")

	done := make(chan struct{})
	<-done
}
