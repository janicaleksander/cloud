package messaging

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"

	"github.com/janicaleksander/cloud/claimservice/application/command"
	"github.com/janicaleksander/cloud/claimservice/domain"
	"github.com/janicaleksander/cloud/common/event"
	"github.com/janicaleksander/cloud/common/rabbitmq"
	"github.com/janicaleksander/cloud/common/rabbitmq/utils"
	"github.com/mehdihadeli/go-mediatr"
)

const exchangeName = "events"
const queueName = "claim-service"

type ClaimEventHandler struct {
	handlers map[string]rabbitmq.HandlerFunc
}

func NewClaimEventHandler() *ClaimEventHandler {
	slog.Info("Creating ClaimEventHandler")
	h := &ClaimEventHandler{
		handlers: make(map[string]rabbitmq.HandlerFunc),
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
		if handler, ok := h.handlers[msg.RoutingKey]; ok {
			handler(&msg)
		} else {
			slog.Error("no handler found for routing key", "routingKey", msg.RoutingKey)
		}
	}
}

func sendChangeStatus(claimID string, status domain.Status) {
	_, err := mediatr.Send[*command.UpdateClaimStatusCommand, *mediatr.Unit](
		context.Background(),
		&command.UpdateClaimStatusCommand{
			ClaimID: claimID,
			Status:  string(status),
		},
	)
	if err != nil {
		slog.Error("failed to change claim status", "status", status, "error", err)
	}
}

func (h *ClaimEventHandler) handlePolicyVerifiedEvent(msg rabbitmq.Delivery) {
	slog.Info("HandlePolicyVerifiedEvent", "routingKey", msg.RoutingKey)
	var e event.PolicyVerifiedEvent
	if err := json.Unmarshal(msg.Body, &e); err != nil {
		slog.Error("failed to unmarshal policy_verified event", "error", err)
		return
	}
	sendChangeStatus(e.ClaimID, domain.VERIFIED)
}

func (h *ClaimEventHandler) handlePolicyDeniedEvent(msg rabbitmq.Delivery) {
	slog.Info("HandlePolicyDeniedEvent", "routingKey", msg.RoutingKey)
	var e event.PolicyDeniedEvent
	if err := json.Unmarshal(msg.Body, &e); err != nil {
		slog.Error("failed to unmarshal policy_denied event", "error", err)
		return
	}
	sendChangeStatus(e.ClaimID, domain.DENIED)
}

func (h *ClaimEventHandler) handlePayoutApprovedEvent(msg rabbitmq.Delivery) {
	slog.Info("HandlePayoutApprovedEvent", "routingKey", msg.RoutingKey)
	var e event.PayoutApprovedEvent
	if err := json.Unmarshal(msg.Body, &e); err != nil {
		slog.Error("failed to unmarshal payout_approved event", "error", err)
		return
	}
	sendChangeStatus(e.ClaimID, domain.APPROVED)
}

func (h *ClaimEventHandler) handlePayoutRejectedEvent(msg rabbitmq.Delivery) {
	slog.Info("HandlePayoutRejectedEvent", "routingKey", msg.RoutingKey)
	var e event.PayoutRejectedEvent
	if err := json.Unmarshal(msg.Body, &e); err != nil {
		slog.Error("failed to unmarshal payout_rejected event", "error", err)
		return
	}
	sendChangeStatus(e.ClaimID, domain.REJECTED)
}
