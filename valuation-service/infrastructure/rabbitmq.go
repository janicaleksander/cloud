package infrastructure

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"time"

	"github.com/janicaleksander/cloud/valuationservice/infrastructure/utils"
	amqp "github.com/rabbitmq/amqp091-go"
)

func NewRabbitMQ() (*amqp.Connection, error) {
	conn, err := amqp.Dial(os.Getenv("AMQP_URL"))
	if err != nil {
		slog.Error("Can't connect to amqp")
		return nil, err
	}
	return conn, nil
}

func Publish(conn *amqp.Connection, exchange string, msg interface{}) error {
	routeKey := utils.NameOfType(msg)

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = ch.PublishWithContext(
		ctx,
		exchange,
		routeKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Timestamp:   time.Now(),
			Body:        b,
		})
	return err
}

func Subscribe[T any](conn *amqp.Connection, exchange string) (<-chan amqp.Delivery, error) {
	var x T
	routeKey := utils.NameOfType(x)

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	q, err := ch.QueueDeclare(
		routeKey,
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(
		q.Name,
		routeKey,
		exchange,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	deliveryChan, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return deliveryChan, nil
}
