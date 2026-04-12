package messaging

import (
	"encoding/json"
	"log/slog"

	"github.com/janicaleksander/cloud/common/event"
	"github.com/janicaleksander/cloud/common/rabbitmq"
	"github.com/janicaleksander/cloud/common/rabbitmq/utils"
	"github.com/janicaleksander/cloud/policyverificationservice/application"
)

const queueName = "policy-verification-service"
const exchangeName = "events"

type PolicyEventHandler struct {
	policyService *application.PolicyService
	handlers      map[string]rabbitmq.HandlerFunc
}

func NewPolicyEventHandler(pS *application.PolicyService) *PolicyEventHandler {
	slog.Info("Initializing PolicyEventHandler")
	p := &PolicyEventHandler{
		policyService: pS,
		handlers:      make(map[string]rabbitmq.HandlerFunc),
	}
	p.registerHandlers()
	return p
}

func (p *PolicyEventHandler) Run(rabbit *rabbitmq.RabbitMQ) {
	slog.Info("Running PolicyEventHandler")
	bindingKeys := make([]string, 0, len(p.handlers))
	for key := range p.handlers {
		bindingKeys = append(bindingKeys, key)
	}
	claimSubmittedChan, err := rabbitmq.SubscribeRaw(rabbit, exchangeName, queueName, bindingKeys...)
	if err != nil {
		slog.Error("Failed to subscribe to RabbitMQ", "error", err)
		return
	}
	go p.dispatch(claimSubmittedChan)
}
func (p *PolicyEventHandler) registerHandlers() {
	slog.Info("Registering event handlers for PolicyEventHandler")
	p.handlers[rabbitmq.RouteKeyToTopicNotation(
		utils.NameOfType(event.ClaimSubmittedEvent{}),
	)] = p.handleClaimSubmittedEvent

}
func (p *PolicyEventHandler) dispatch(msgs rabbitmq.MsgChan) {
	for msg := range msgs {
		if handler, ok := p.handlers[msg.RoutingKey]; ok {
			handler(&msg)
		} else {
			slog.Error("Unknown routing key", "routingKey", msg.RoutingKey)
		}
	}

}

func (p *PolicyEventHandler) handleClaimSubmittedEvent(msg rabbitmq.Delivery) {
	slog.Info("HandleClaimSubmittedEvent: ", "routingKey", msg.RoutingKey)
	var claimSubmittedEvent event.ClaimSubmittedEvent
	err := json.Unmarshal(msg.Body, &claimSubmittedEvent)
	if err != nil {
		slog.Info("Failed to unmarshal ClaimSubmittedEvent", "error", err)
		return
	}
	p.policyService.CheckUserPolicy(
		claimSubmittedEvent.ClaimID,
		claimSubmittedEvent.UserID,
		claimSubmittedEvent.VIN,
		claimSubmittedEvent.AccidentDate,
		claimSubmittedEvent.StorageURL,
	)
}
