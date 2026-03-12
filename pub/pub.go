package pub

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/utils"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	ID      uuid.UUID
	Channel *amqp.Channel
}

func NewPublisher(conn *amqp.Connection) (*Publisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	p := &Publisher{
		ID:      uuid.New(),
		Channel: ch,
	}
	defer func() {
		slog.Info("Creating publisher with ID: " + p.ID.String())
	}()
	return p, nil

}
func (p *Publisher) Publish(ctx context.Context, msg interface{}) error {
	slog.Info(fmt.Sprintf("Publishing message: %v, type: %s, by publisher %s", msg, utils.GetTypeName(msg), p.ID))
	bytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = p.Channel.PublishWithContext(ctx,
		"",
		utils.GetTypeName(msg),
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bytes,
		})
	return err
}
