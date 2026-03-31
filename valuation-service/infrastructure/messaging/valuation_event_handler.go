package messaging

import (
	"encoding/json"

	"github.com/janicaleksander/cloud/common/event"
	"github.com/janicaleksander/cloud/common/rabbitmq"
	"github.com/janicaleksander/cloud/common/rabbitmq/utils"
	"github.com/janicaleksander/cloud/valuationservice/application"
)

type ValuationEventHandler struct {
	valuationService *application.ValuationService
	handlers         map[string]rabbitmq.HandlerFunc
}

func NewValuationEventHandler(vS *application.ValuationService) *ValuationEventHandler {
	v := &ValuationEventHandler{
		valuationService: vS,
		handlers:         make(map[string]rabbitmq.HandlerFunc),
	}
	v.registerHandlers()
	return v
}

func (v *ValuationEventHandler) Run(rabbit *rabbitmq.RabbitMQ) {
	bindingKeys := make([]string, 0, len(v.handlers))
	for key := range v.handlers {
		bindingKeys = append(bindingKeys, key)
	}
	claimSubmittedChan, err := rabbitmq.SubscribeRaw(rabbit, "events", "valuation-service", bindingKeys...)
	if err != nil {
		//TODO logs
		return
	}
	go v.dispatch(claimSubmittedChan)
}

func (v *ValuationEventHandler) registerHandlers() {
	v.handlers[rabbitmq.RouteKeyToTopicNotation(
		utils.NameOfType(event.PolicyVerifiedEvent{}),
	)] = v.handlePolicyVerifiedEvent
}

func (v *ValuationEventHandler) handlePolicyVerifiedEvent(msg rabbitmq.Delivery) {
	var policyVerifiedEvent event.PolicyVerifiedEvent
	err := json.Unmarshal(msg.Body, &policyVerifiedEvent)
	if err != nil {
		//TODO logs
		return
	}
	v.valuationService.CalculateValuation(
		policyVerifiedEvent.StorageURL,
		policyVerifiedEvent.ClaimID,
	)

}

func (v *ValuationEventHandler) dispatch(msgs rabbitmq.MsgChan) {
	for msg := range msgs {
		if handler, ok := v.handlers[msg.RoutingKey]; ok {
			handler(&msg)
		} else {
			//TODO logs
		}
	}

}
