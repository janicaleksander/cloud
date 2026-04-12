package messaging

import (
	"encoding/json"
	"log/slog"

	"github.com/janicaleksander/cloud/common/event"
	"github.com/janicaleksander/cloud/common/rabbitmq"
	"github.com/janicaleksander/cloud/common/rabbitmq/utils"
	"github.com/janicaleksander/cloud/decisionservice/application"
)

const queueName = "decision-services"
const exchangeName = "events"

type DecisionEventHandler struct {
	decisionService *application.DecisionService
	handlers        map[string]rabbitmq.HandlerFunc
}

func NewDecisionEventHandler(decisionService *application.DecisionService) *DecisionEventHandler {
	slog.Info("Creating DecisionEventHandler")
	h := &DecisionEventHandler{
		decisionService: decisionService,
		handlers:        make(map[string]rabbitmq.HandlerFunc),
	}
	h.registerHandlers()
	return h
}

func (dH *DecisionEventHandler) registerHandlers() {
	slog.Info("Registering handlers for DecisionEventHandler")
	dH.handlers[rabbitmq.RouteKeyToTopicNotation(utils.NameOfType(event.ValuationCalculatedEvent{}))] = dH.handleValuationCalculated

}

func (dH *DecisionEventHandler) Run(rabbit *rabbitmq.RabbitMQ) error {
	slog.Info("Running DecisionEventHandler")
	bindingKeys := make([]string, 0, len(dH.handlers))
	for k, _ := range dH.handlers {
		bindingKeys = append(bindingKeys, k)
	}
	msgs, err := rabbitmq.SubscribeRaw(rabbit, exchangeName, queueName, bindingKeys...)
	if err != nil {
		return err
	}
	go dH.dispatch(msgs)
	return nil
}

func (dH *DecisionEventHandler) dispatch(msgChan rabbitmq.MsgChan) {
	for msg := range msgChan {
		if handler, ok := dH.handlers[msg.RoutingKey]; ok {
			handler(&msg)

		} else {
			slog.Error("Unknown routing key", "routingKey", msg.RoutingKey)
		}
	}

}
func (dH *DecisionEventHandler) handleValuationCalculated(msg rabbitmq.Delivery) {
	slog.Info("Handling ValuationCalculatedEvent", "routingKey", msg.RoutingKey)
	var e event.ValuationCalculatedEvent
	err := json.Unmarshal(msg.Body, &e)
	if err != nil {
		slog.Error("Failed to unmarshal ValuationCalculatedEvent", "error", err)
		return
	}
	_, err = dH.decisionService.PrepareDecision(e.ClaimID, e.PayoutAmount)
	if err != nil {
		slog.Error("Failed to prepare decision", "error", err)
	}
}
