package messaging

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/janicaleksander/cloud/claimservice/application"
	"github.com/janicaleksander/cloud/claimservice/domain"
	"github.com/janicaleksander/cloud/common/event"
	"github.com/janicaleksander/cloud/common/rabbitmq"
	"github.com/janicaleksander/cloud/common/rabbitmq/utils"
)

type ClaimEventHandler struct {
	claimService *application.ClaimService
	handlers     map[string]rabbitmq.HandlerFunc
}

func NewClaimEventHandler(claimService *application.ClaimService) *ClaimEventHandler {
	h := &ClaimEventHandler{
		claimService: claimService,
		handlers:     make(map[string]rabbitmq.HandlerFunc),
	}
	h.registerHandlers()
	return h
}

func (h *ClaimEventHandler) registerHandlers() {
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
	bindingKeys := make([]string, 0, len(h.handlers))
	for key := range h.handlers {
		bindingKeys = append(bindingKeys, key)
	}

	msgs, err := rabbitmq.SubscribeRaw(rabbit, "events", "claim-service", bindingKeys...)
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
			log.Println("err")
		}
	}
}

func (h *ClaimEventHandler) handlePolicyVerifiedEvent(msg rabbitmq.Delivery) {
	log.Println("HandlePolicyVerifiedEvent")
	var e event.PolicyVerifiedEvent
	err := json.Unmarshal(msg.Body, &e)
	if err != nil {
		log.Printf("failed to unmarshal policy_verified event: %v", err)
		//TODO log this
		return
	}
	err = h.claimService.ChangeClaimStatus(e.ClaimID, domain.VERIFIED)
	if err != nil {
		log.Printf("failed to change claim status to VERIFIED: %v", err)
		//TODO log
	}
}

func (h *ClaimEventHandler) handlePolicyDeniedEvent(msg rabbitmq.Delivery) {
	log.Printf("HandlePolicyDeniedEvent: ")
	var e event.PolicyDeniedEvent
	err := json.Unmarshal(msg.Body, &e)
	if err != nil {
		log.Printf("failed to unmarshal policy_denied event: %v", err)
		//TODO log this
		return
	}
	err = h.claimService.ChangeClaimStatus(e.ClaimID, domain.DENIED)
	if err != nil {
		log.Printf("failed to change claim status to DENIED: %v", err)
		//TODO log
	}

}

func (h *ClaimEventHandler) handlePayoutApprovedEvent(msg rabbitmq.Delivery) {
	log.Printf("HandlePayoutApprovedEvent: ")
	var e event.PayoutApprovedEvent
	err := json.Unmarshal(msg.Body, &e)
	if err != nil {
		log.Printf("failed to unmarshal payout_approved event: %v", err)
		//TODO log this
		return
	}
	err = h.claimService.ChangeClaimStatus(e.ClaimID, domain.APPROVED)
	if err != nil {
		log.Printf("failed to change claim status to APPROVED: %v", err)
		//TODO log
	}
}

func (h *ClaimEventHandler) handlePayoutRejectedEvent(msg rabbitmq.Delivery) {
	log.Printf("HandlePayoutRejectedEvent: ")
	var e event.PayoutRejectedEvent
	err := json.Unmarshal(msg.Body, &e)
	if err != nil {
		log.Printf("failed to unmarshal payout_rejected event: %v", err)
		//TODO log this
		return
	}
	err = h.claimService.ChangeClaimStatus(e.ClaimID, domain.REJECTED)
	if err != nil {
		log.Printf("failed to change claim status to REJECTED: %v", err)
		//TODO log
	}

}
