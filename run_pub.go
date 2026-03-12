package main

import (
	"context"
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

const (
	type1Delay       = 2 * time.Second
	publishTimeout   = 5 * time.Second
	randomDelayRange = 5
	type1Count       = 3
)

func startPublisher(p *pub.Publisher, generator func() any, delay func() time.Duration) {
	go func() {
		for {
			ctx, cancel := context.WithTimeout(context.Background(), publishTimeout)
			if err := p.Publish(ctx, generator()); err != nil {
				slog.Error("Failed to publish: " + err.Error())
			}
			cancel()
			time.Sleep(delay())
		}
	}()
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
	exchange := os.Getenv("EXCHANGE_NAME")

	publishers := make([]*pub.Publisher, 5)
	for i := range publishers {
		publishers[i], err = pub.NewPublisher(conn, exchange)
		if err != nil {
			slog.Error("Spawning publisher error ", err.Error())
			return
		}
		defer publishers[i].Channel.Close()
	}

	for i := 0; i < type1Count; i++ {
		startPublisher(publishers[i], func() any { return event.NewType1Event() }, utils.Delay(type1Delay))
	}

	startPublisher(publishers[type1Count], func() any { return event.NewType2Event() }, utils.Delay(time.Second*time.Duration(rand.Intn(randomDelayRange)+1)))
	startPublisher(publishers[type1Count+1], func() any { return event.NewType3Event() }, utils.Delay(time.Second*time.Duration(rand.Intn(randomDelayRange)+1)))

	done := make(chan struct{})
	<-done
}
