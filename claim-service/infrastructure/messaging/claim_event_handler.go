package messaging

import (
	"encoding/json"
	"log"

	"github.com/janicaleksander/cloud/claimservice/application"
	"github.com/janicaleksander/cloud/claimservice/domain"
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
		log.Println("HandlePolicyVerifiedEvent")
		var e event.PolicyVerifiedEvent
		err := json.Unmarshal(msg.Body, &e)
		if err != nil {
			log.Printf("failed to unmarshal policy_verified event: %v", err)
			//TODO log this
		}
		err = h.claimService.ChangeClaimStatus(e.ClaimID, domain.VERIFIED)
		if err != nil {
			log.Printf("failed to change claim status to VERIFIED: %v", err)
			//TODO log
		}
	}
}

func (h *ClaimEventHandler) handlePolicyDeniedEvent(msgs rabbitmq.MsgChan) {
	for msg := range msgs {
		log.Printf("HandlePolicyDeniedEvent: ")
		var e event.PolicyDeniedEvent
		err := json.Unmarshal(msg.Body, &e)
		if err != nil {
			log.Printf("failed to unmarshal policy_denied event: %v", err)
			//TODO log this
		}
		err = h.claimService.ChangeClaimStatus(e.ClaimID, domain.DENIED)
		if err != nil {
			log.Printf("failed to change claim status to DENIED: %v", err)
			//TODO log
		}

	}
}

func (h *ClaimEventHandler) handlePayoutApprovedEvent(msgs rabbitmq.MsgChan) {
	for msg := range msgs {
		log.Printf("HandlePayoutApprovedEvent: ")
		var e event.PayoutApprovedEvent
		err := json.Unmarshal(msg.Body, &e)
		if err != nil {
			log.Printf("failed to unmarshal payout_approved event: %v", err)
			//TODO log this
		}
		err = h.claimService.ChangeClaimStatus(e.ClaimID, domain.APPROVED)
		if err != nil {
			log.Printf("failed to change claim status to APPROVED: %v", err)
			//TODO log
		}

	}
}

func (h *ClaimEventHandler) handlePayoutRejectedEvent(msgs rabbitmq.MsgChan) {
	for msg := range msgs {
		log.Printf("HandlePayoutRejectedEvent: ")
		var e event.PayoutRejectedEvent
		err := json.Unmarshal(msg.Body, &e)
		if err != nil {
			log.Printf("failed to unmarshal payout_rejected event: %v", err)
			//TODO log this
		}
		err = h.claimService.ChangeClaimStatus(e.ClaimID, domain.REJECTED)
		if err != nil {
			log.Printf("failed to change claim status to REJECTED: %v", err)
			//TODO log
		}

	}
}
