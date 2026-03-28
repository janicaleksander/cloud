package messaging

import (
	"log"

	"github.com/janicaleksander/cloud/claimservice/application"
	"github.com/janicaleksander/cloud/common/event"
	"github.com/janicaleksander/cloud/common/rabbitmq"
)

type ClaimEventHandler struct {
	claimService *application.ClaimService
}

func NewClaimEventHandler(claimService *application.ClaimService) *ClaimEventHandler {
	return &ClaimEventHandler{claimService: claimService}
}

func (h *ClaimEventHandler) Run(rabbit *rabbitmq.RabbitMQ) {
	policyVerifiedChan, err := rabbitmq.Subscribe[event.PolicyVerifiedEvent](rabbit, "events")
	if err != nil {
		log.Printf("failed to subscribe to policy_verified: %v", err)
	}

	policyDeniedChan, err := rabbitmq.Subscribe[event.PolicyDeniedEvent](rabbit, "events")
	if err != nil {
		log.Printf("failed to subscribe to policy_denied: %v", err)
	}

	payoutApprovedChan, err := rabbitmq.Subscribe[event.PayoutApprovedEvent](rabbit, "events")
	if err != nil {
		log.Printf("failed to subscribe to payout_approved: %v", err)
	}

	payoutRejectedChan, err := rabbitmq.Subscribe[event.PayoutRejectedEvent](rabbit, "events")
	if err != nil {
		log.Printf("failed to subscribe to payout_rejected: %v", err)
	}

	go h.handlePolicyVerifiedEvent(policyVerifiedChan)
	go h.handlePolicyDeniedEvent(policyDeniedChan)
	go h.handlePayoutApprovedEvent(payoutApprovedChan)
	go h.handlePayoutRejectedEvent(payoutRejectedChan)
}

func (h *ClaimEventHandler) handlePolicyVerifiedEvent(msgs rabbitmq.MsgChan) {
	for msg := range msgs {
		log.Printf("HandlePolicyVerifiedEvent: %+v", msg)
		//h.claimService.UpdateClaim()
	}
}

func (h *ClaimEventHandler) handlePolicyDeniedEvent(msgs rabbitmq.MsgChan) {
	for msg := range msgs {
		log.Printf("HandlePolicyDeniedEvent: %+v", msg)
		//h.claimService.UpdateClaim()

	}
}

func (h *ClaimEventHandler) handlePayoutApprovedEvent(msgs rabbitmq.MsgChan) {
	for msg := range msgs {
		log.Printf("HandlePayoutApprovedEvent: %+v", msg)
		//h.claimService.UpdateClaim()

	}
}

func (h *ClaimEventHandler) handlePayoutRejectedEvent(msgs rabbitmq.MsgChan) {
	for msg := range msgs {
		log.Printf("HandlePayoutRejectedEvent: %+v", msg)
		//h.claimService.UpdateClaim()

	}
}
