package messaging

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/janicaleksander/cloud/common/event"
	"github.com/janicaleksander/cloud/common/rabbitmq"
	"github.com/janicaleksander/cloud/common/rabbitmq/utils"
	"github.com/janicaleksander/cloud/valuationservice/application/command"
	"github.com/mehdihadeli/go-mediatr"
)

const queueName = "valuation-service"
const exchangeName = "events"

type ValuationEventHandler struct {
	handlers map[string]rabbitmq.HandlerFunc
}

func NewValuationEventHandler() *ValuationEventHandler {
	slog.Info("Creating ValuationEventHandler")
	v := &ValuationEventHandler{
		handlers: make(map[string]rabbitmq.HandlerFunc),
	}
	v.registerHandlers()
	return v
}

func (v *ValuationEventHandler) Run(rabbit *rabbitmq.RabbitMQ) {
	slog.Info("Running ValuationEventHandler")
	bindingKeys := make([]string, 0, len(v.handlers))
	for key := range v.handlers {
		bindingKeys = append(bindingKeys, key)
	}
	claimSubmittedChan, err := rabbitmq.SubscribeRaw(rabbit, exchangeName, queueName, bindingKeys...)
	if err != nil {
		slog.Error("Failed to subscribe to RabbitMQ", "error", err.Error())
		return
	}
	go v.dispatch(claimSubmittedChan)
}

func (v *ValuationEventHandler) registerHandlers() {
	slog.Info("Registering event handlers for ValuationEventHandler")
	v.handlers[rabbitmq.RouteKeyToTopicNotation(
		utils.NameOfType(event.PolicyVerifiedEvent{}),
	)] = v.handlePolicyVerifiedEvent
}

func (v *ValuationEventHandler) handlePolicyVerifiedEvent(msg rabbitmq.Delivery) {
	slog.Info("HandlePolicyVerifiedEvent", "routingKey", msg.RoutingKey)
	var policyVerifiedEvent event.PolicyVerifiedEvent
	err := json.Unmarshal(msg.Body, &policyVerifiedEvent)
	if err != nil {
		slog.Error("failed to unmarshal PolicyVerifiedEvent", "error", err.Error())
		return
	}
	cmd := &command.CalculateValuationCommand{
		ClaimID: policyVerifiedEvent.ClaimID,
		Urls:    policyVerifiedEvent.StorageURL,
	}
	_, err = mediatr.Send[*command.CalculateValuationCommand, *mediatr.Unit](context.Background(), cmd)
	if err != nil {
		slog.Error("failed to send CalculateValuationCommand", "error", err.Error())
		return
	}
}

func (v *ValuationEventHandler) dispatch(msgs rabbitmq.MsgChan) {
	for msg := range msgs {
		if handler, ok := v.handlers[msg.RoutingKey]; ok {
			handler(&msg)
		} else {
			slog.Error("Unknown routing key", "routingKey", msg.RoutingKey)
		}
	}

}
