package rabbitmq

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/janicaleksander/cloud/common/rabbitmq/utils"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	rabbit *RabbitMQ
}

type RabbitMQ struct {
	conn *amqp.Connection
}
type HandlerFunc func(msgs Delivery)

type MsgChan <-chan amqp.Delivery
type Delivery *amqp.Delivery

func NewPublisher(r *RabbitMQ) *Publisher {
	slog.Info("Creating RabbitMQ publisher")
	return &Publisher{
		rabbit: r,
	}
}

func NewRabbitMQ() (*RabbitMQ, error) {
	slog.Info("Connecting to RabbitMQ")
	conn, err := amqp.Dial(os.Getenv("AMQP_URL"))
	if err != nil {
		slog.Error("Can't connect to amqp")
		return nil, err
	}
	return &RabbitMQ{conn: conn}, nil
}

func (p *Publisher) Publish(exchange string, msg interface{}) error {
	slog.Info("Publishing message to RabbitMQ", "exchange", exchange, "type", utils.NameOfType(msg))
	routeKey := RouteKeyToTopicNotation(utils.NameOfType(msg))
	ch, err := p.rabbit.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	err = ch.ExchangeDeclare(
		exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
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

func RouteKeyToTopicNotation(routeKey string) string {
	if strings.HasSuffix(routeKey, "Event") {
		routeKey = routeKey[:len(routeKey)-5]
	}

	var result []rune
	for i, r := range routeKey {
		if i > 0 && unicode.IsUpper(r) &&
			(i+1 < len(routeKey) && unicode.IsLower(rune(routeKey[i+1]))) {
			result = append(result, '.')
		}
		result = append(result, unicode.ToLower(r))
	}

	return string(result)
}
func Subscribe[T any](rabbitmq *RabbitMQ, exchange string, qName string) (<-chan amqp.Delivery, error) {
	slog.Info("Subscribing to RabbitMQ", "exchange", exchange)
	var x T
	routeKey := RouteKeyToTopicNotation(utils.NameOfType(x))

	ch, err := rabbitmq.conn.Channel()
	if err != nil {
		return nil, err
	}
	err = ch.ExchangeDeclare(
		exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	q, err := ch.QueueDeclare(
		qName,
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

func SubscribeRaw(rabbitmq *RabbitMQ, exchange string, qName string, bindingKeys ...string) (<-chan amqp.Delivery, error) {
	slog.Info("Subscribing to RabbitMQ", "exchange", exchange)
	ch, err := rabbitmq.conn.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(exchange, "topic", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(qName, false, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	if len(bindingKeys) == 0 {
		bindingKeys = []string{"#"}
	}

	for _, key := range bindingKeys {
		err = ch.QueueBind(q.Name, key, exchange, false, nil)
		if err != nil {
			return nil, err
		}
	}

	return ch.Consume(q.Name, "", true, false, false, false, nil)
}
