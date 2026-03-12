package pub

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/utils"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	ID      uuid.UUID
	channel *amqp.Channel
}

func NewPublisher(conn *amqp.Connection) (*Publisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &Publisher{
		ID:      uuid.New(),
		channel: ch,
	}, nil

}
func (p *Publisher) Publish(msg interface{}) {
	slog.Info(fmt.Sprintf("Publishing message: %v, type: %s, by publisher %s", msg, utils.GetTypeName(msg), p.ID))
	bytes, err := json.Marshal(msg)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	err = p.channel.Publish(
		"",
		utils.GetTypeName(msg),
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bytes,
		})
	if err != nil {
		slog.Error(err.Error())
	}
}
