package messaging

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"

	"github.com/janicaleksander/cloud/claimservice/application"
	"github.com/janicaleksander/cloud/claimservice/domain"
	"github.com/janicaleksander/cloud/common/event"
	"github.com/janicaleksander/cloud/common/rabbitmq"
	"github.com/janicaleksander/cloud/common/rabbitmq/utils"
)

const exchangeName = "events"
const queueName = "claim-service"

type ClaimEventHandler struct {
	claimService *application.ClaimService
	handlers     map[string]rabbitmq.HandlerFunc
}

func NewClaimEventHandler(claimService *application.ClaimService) *ClaimEventHandler {
	slog.Info("Creating ClaimEventHandler")
	h := &ClaimEventHandler{
		claimService: claimService,
		handlers:     make(map[string]rabbitmq.HandlerFunc),
	}
	h.registerHandlers()
	return h
}

func (h *ClaimEventHandler) registerHandlers() {
	slog.Info("Registering event handlers for ClaimEventHandler")
	h.handlers[rabbitmq.RouteKeyToTopicNotation(
		utils.NameOfType(event.PolicyVerifiedEvent{}),
	)] = h.handlePolicyVerifiedEvent

	h.handlers[rabbitmq.RouteKeyToTopicNotation(
		utils.NameOfType(event.PolicyDeniedEvent{}),
	)] = h.handlePolicyDeniedEvent

	h.handlers[rabbitmq.RouteKeyToTopicNotation(
		utils.NameOfType(event.PayoutApprovedEvent{}),
	)] = h.handlePayoutApprovedEvent

	h.handlers[rabbitmq.RouteKeyToTopicNotation(
		utils.NameOfType(event.PayoutRejectedEvent{}),
	)] = h.handlePayoutRejectedEvent
}

func (h *ClaimEventHandler) Run(rabbit *rabbitmq.RabbitMQ) {
	slog.Info("Running ClaimEventHandler")
	bindingKeys := make([]string, 0, len(h.handlers))
	for key := range h.handlers {
		bindingKeys = append(bindingKeys, key)
	}

	msgs, err := rabbitmq.SubscribeRaw(rabbit, exchangeName, queueName, bindingKeys...)
	if err != nil {
		log.Fatal(err)
	}

	go h.dispatch(msgs)
}
func (h *ClaimEventHandler) dispatch(msgs rabbitmq.MsgChan) {
	for msg := range msgs {
		fmt.Println(msg.RoutingKey)
		if handler, ok := h.handlers[msg.RoutingKey]; ok {
			handler(&msg)
		} else {
			slog.Error("no handler found for routing key", "routingKey", msg.RoutingKey)
		}
	}
}

func (h *ClaimEventHandler) handlePolicyVerifiedEvent(msg rabbitmq.Delivery) {
	slog.Info("HandlePolicyVerifiedEvent: ", "routingKey", msg.RoutingKey)
	var e event.PolicyVerifiedEvent
	err := json.Unmarshal(msg.Body, &e)
	if err != nil {
		slog.Error("failed to unmarshal policy_verified event", "error", err)
		return
	}
	err = h.claimService.ChangeClaimStatus(e.ClaimID, domain.VERIFIED)
	if err != nil {
		slog.Error("failed to change claim status to VERIFIED", "error", err)
	}
}

func (h *ClaimEventHandler) handlePolicyDeniedEvent(msg rabbitmq.Delivery) {
	slog.Info("HandlePolicyDeniedEvent: ", "routingKey", msg.RoutingKey)
	var e event.PolicyDeniedEvent
	err := json.Unmarshal(msg.Body, &e)
	if err != nil {
		slog.Error("failed to unmarshal policy_denied event", "error", err)
		return
	}
	err = h.claimService.ChangeClaimStatus(e.ClaimID, domain.DENIED)
	if err != nil {
		slog.Error("failed to change claim status to DENIED", "error", err)
	}

}

func (h *ClaimEventHandler) handlePayoutApprovedEvent(msg rabbitmq.Delivery) {
	slog.Info("HandlePayoutApprovedEvent: ", "routingKey", msg.RoutingKey)
	var e event.PayoutApprovedEvent
	err := json.Unmarshal(msg.Body, &e)
	if err != nil {
		slog.Error("failed to unmarshal payout_approved event", "error", err)
		return
	}
	err = h.claimService.ChangeClaimStatus(e.ClaimID, domain.APPROVED)
	if err != nil {
		slog.Error("failed to change claim status to APPROVED", "error", err)
	}
}

func (h *ClaimEventHandler) handlePayoutRejectedEvent(msg rabbitmq.Delivery) {
	slog.Info("HandlePayoutRejectedEvent: ", "routingKey", msg.RoutingKey)
	var e event.PayoutRejectedEvent
	err := json.Unmarshal(msg.Body, &e)
	if err != nil {
		slog.Error("failed to unmarshal payout_rejected event", "error", err)
		return
	}
	err = h.claimService.ChangeClaimStatus(e.ClaimID, domain.REJECTED)
	if err != nil {
		slog.Error("failed to change claim status to REJECTED", "error", err)
	}

}
