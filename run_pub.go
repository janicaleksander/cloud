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

/*
1.	3 publisherów generujących zdarzenia typu 1, publikujących je w takich samych odstępach czasowych
2.	1 publishera generującego zdarzenia typu 2, publikujący je w losowym odstępie czasowym
3.	1 publishera generującego zdarzenia typu 3, publikujący je w losowym odstępie czasowym.
*/
const (
	type1Delay       = 2 * time.Second
	publishTimeout   = 5 * time.Second
	randomDelayRange = 5
)

func publishWithRandomDelay(p *pub.Publisher, eventGenerator func() interface{}, eventType string) {
	go func() {
		for {
			ctx, cancel := context.WithTimeout(context.Background(), publishTimeout)
			err := p.Publish(ctx, eventGenerator())
			if err != nil {
				slog.Error("Failed to publish " + eventType + ": " + err.Error())
			}
			cancel()
			time.Sleep(time.Second * time.Duration(rand.Intn(randomDelayRange)+1))
		}
	}()
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

	p11, err := pub.NewPublisher(conn)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer p11.Channel.Close()

	p21, err := pub.NewPublisher(conn)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer p21.Channel.Close()

	p31, err := pub.NewPublisher(conn)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer p31.Channel.Close()
	type1Publishers := []*pub.Publisher{p11, p21, p31}
	go func() {
		for {
			ctx, cancel := context.WithTimeout(context.Background(), publishTimeout)
			for _, p := range type1Publishers {
				if err := p.Publish(ctx, event.NewType1Event()); err != nil {
					slog.Error("Failed to publish Type1Event: " + err.Error())
				}
			}
			cancel()
			time.Sleep(type1Delay)
		}
	}()

	p12, err := pub.NewPublisher(conn)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer p12.Channel.Close()
	publishWithRandomDelay(p12, func() interface{} { return event.NewType2Event() }, utils.GetTypeName(event.NewType2Event()))

	p13, err := pub.NewPublisher(conn)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer p13.Channel.Close()
	publishWithRandomDelay(p13, func() interface{} { return event.NewType3Event() }, utils.GetTypeName(event.NewType3Event()))

	<-done
}
