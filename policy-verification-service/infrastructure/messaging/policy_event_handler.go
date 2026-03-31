package messaging

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/janicaleksander/cloud/common/event"
	"github.com/janicaleksander/cloud/common/rabbitmq"
	"github.com/janicaleksander/cloud/common/rabbitmq/utils"
	"github.com/janicaleksander/cloud/policyverificationservice/application"
)

type PolicyEventHandler struct {
	policyService *application.PolicyService
	handlers      map[string]rabbitmq.HandlerFunc
}

func NewPolicyEventHandler(pS *application.PolicyService) *PolicyEventHandler {
	p := &PolicyEventHandler{
		policyService: pS,
		handlers:      make(map[string]rabbitmq.HandlerFunc),
	}
	p.registerHandlers()
	return p
}

func (p *PolicyEventHandler) Run(rabbit *rabbitmq.RabbitMQ) {
	bindingKeys := make([]string, 0, len(p.handlers))
	for key := range p.handlers {
		bindingKeys = append(bindingKeys, key)
	}
	claimSubmittedChan, err := rabbitmq.SubscribeRaw(rabbit, "events", "policy-verification-service", bindingKeys...)
	if err != nil {
		//TODO logs
		return
	}
	go p.dispatch(claimSubmittedChan)
}
func (p *PolicyEventHandler) registerHandlers() {
	p.handlers[rabbitmq.RouteKeyToTopicNotation(
		utils.NameOfType(event.ClaimSubmittedEvent{}),
	)] = p.handleClaimSubmittedEvent

}
func (p *PolicyEventHandler) dispatch(msgs rabbitmq.MsgChan) {
	for msg := range msgs {
		fmt.Println(msg.RoutingKey)
		if handler, ok := p.handlers[msg.RoutingKey]; ok {
			handler(&msg)
		} else {
			log.Println("err")
		}
	}

}

func (p *PolicyEventHandler) handleClaimSubmittedEvent(msg rabbitmq.Delivery) {
	var claimSubmittedEvent event.ClaimSubmittedEvent
	log.Println("Received ClaimSubmittedEvent:", string(msg.Body))
	err := json.Unmarshal(msg.Body, &claimSubmittedEvent)
	if err != nil {
		log.Println("Error unmarshalling ClaimSubmittedEvent:", err)
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
