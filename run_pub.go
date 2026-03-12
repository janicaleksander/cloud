package main

import (
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
	type1Delay = 2 * time.Second
)

// TOOD repair connection leak, close the channels and connections
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
	go func() {
		for {
			p11.Publish(event.NewType1Event())
			p21.Publish(event.NewType1Event())
			p31.Publish(event.NewType1Event())
			time.Sleep(type1Delay)
		}
	}()
	p12, err := pub.NewPublisher(conn)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer p12.Channel.Close()
	go func() {
		for {
			p12.Publish(event.NewType2Event())
			time.Sleep(time.Second * time.Duration(rand.Intn(5)+1))
		}
	}()
	p13, err := pub.NewPublisher(conn)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer p13.Channel.Close()
	go func() {
		for {
			p13.Publish(event.NewType3Event())
			time.Sleep(time.Second * time.Duration(rand.Intn(5)+1))
		}
	}()

	<-done
}
