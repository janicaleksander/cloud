package messaging

import (
	"encoding/json"
	"log"

	"github.com/janicaleksander/cloud/common/event"
	"github.com/janicaleksander/cloud/common/rabbitmq"
	"github.com/janicaleksander/cloud/policyverificationservice/application"
)

type PolicyEventHandler struct {
	policyService *application.PolicyService
}

func NewPolicyEventHandler(pS *application.PolicyService) *PolicyEventHandler {
	return &PolicyEventHandler{policyService: pS}
}

func (p *PolicyEventHandler) Run(rabbit *rabbitmq.RabbitMQ) {
	claimSubmittedChan, err := rabbitmq.Subscribe[event.ClaimSubmittedEvent](rabbit, "events", "policy-verification-service")
	if err != nil {
		//TODO logs
	}
	go p.handleClaimSubmittedEvent(claimSubmittedChan)
}

func (p *PolicyEventHandler) handleClaimSubmittedEvent(msgs rabbitmq.MsgChan) {
	var claimSubmittedEvent event.ClaimSubmittedEvent
	for msg := range msgs {
		log.Println("Received ClaimSubmittedEvent:", string(msg.Body))
		err := json.Unmarshal(msg.Body, &claimSubmittedEvent)
		if err != nil {
			log.Println("Error unmarshalling ClaimSubmittedEvent:", err)
			continue
		}
		p.policyService.CheckUserPolicy(
			claimSubmittedEvent.ClaimID,
			claimSubmittedEvent.UserID,
			claimSubmittedEvent.VIN,
			claimSubmittedEvent.AccidentDate,
			claimSubmittedEvent.StorageURL,
		)

	}
}
