package main

import (
	"context"
	"log/slog"
	"math/rand"
	"os"
	"time"

	"github.com/janicaleksander/cloud/event"
	"github.com/janicaleksander/cloud/pub"
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

func createPublisher(conn *amqp.Connection) (*pub.Publisher, error) {
	p, err := pub.NewPublisher(conn)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func publishWithContext(p *pub.Publisher, msg interface{}, eventType string) {
	ctx, cancel := context.WithTimeout(context.Background(), publishTimeout)
	defer cancel()

	err := p.Publish(ctx, msg)
	if err != nil {
		slog.Error("Failed to publish " + eventType + ": " + err.Error())
	}
}

func runType1Publishers(publishers []*pub.Publisher) {
	go func() {
		for {
			for _, p := range publishers {
				publishWithContext(p, event.NewType1Event(), "Type1Event")
			}
			time.Sleep(type1Delay)
		}
	}()
}

func runRandomPublisher(p *pub.Publisher, eventGenerator func() interface{}, eventType string) {
	go func() {
		for {
			publishWithContext(p, eventGenerator(), eventType)
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
		slog.Error("Can't connect to amqp")
		return
	}
	defer conn.Close()

	// Create Type1Event publishers (3 publishers)
	type1Publishers := make([]*pub.Publisher, 3)
	for i := 0; i < 3; i++ {
		p, err := createPublisher(conn)
		if err != nil {
			slog.Error("Failed to create Type1Event publisher: " + err.Error())
			return
		}
		defer p.Channel.Close()
		type1Publishers[i] = p
	}
	runType1Publishers(type1Publishers)

	// Create and run Type2Event publisher
	p12, err := createPublisher(conn)
	if err != nil {
		slog.Error("Failed to create Type2Event publisher: " + err.Error())
		return
	}
	defer p12.Channel.Close()
	runRandomPublisher(p12, func() interface{} { return event.NewType2Event() }, "Type2Event")

	// Create and run Type3Event publisher
	p13, err := createPublisher(conn)
	if err != nil {
		slog.Error("Failed to create Type3Event publisher: " + err.Error())
		return
	}
	defer p13.Channel.Close()
	runRandomPublisher(p13, func() interface{} { return event.NewType3Event() }, "Type3Event")

	done := make(chan struct{})
	<-done
}
