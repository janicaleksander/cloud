package messaging

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/janicaleksander/cloud/common/event"
	"github.com/janicaleksander/cloud/common/rabbitmq"
	"github.com/janicaleksander/cloud/common/rabbitmq/utils"
	"github.com/janicaleksander/cloud/decisionservice/application"
)

type DecisionEventHandler struct {
	decisionService *application.DecisionService
	handlers        map[string]rabbitmq.HandlerFunc
}

func NewDecisionEventHandler(decisionService *application.DecisionService) *DecisionEventHandler {
	h := &DecisionEventHandler{
		decisionService: decisionService,
		handlers:        make(map[string]rabbitmq.HandlerFunc),
	}
	h.registerHandlers()
	return h
}

func (dH *DecisionEventHandler) registerHandlers() {
	dH.handlers[rabbitmq.RouteKeyToTopicNotation(utils.NameOfType(event.ValuationCalculatedEvent{}))] = dH.handleValuationCalculated

}

func (dH *DecisionEventHandler) Run(rabbit *rabbitmq.RabbitMQ) error {
	fmt.Println("Running DecisionEventHandler")
	bindingKeys := make([]string, 0, len(dH.handlers))
	for k, _ := range dH.handlers {
		bindingKeys = append(bindingKeys, k)
	}
	msgs, err := rabbitmq.SubscribeRaw(rabbit, "events", "decision-services", bindingKeys...)
	if err != nil {
		return err
	}
	go dH.dispatch(msgs)
	return nil
}

func (dH *DecisionEventHandler) dispatch(msgChan rabbitmq.MsgChan) {
	for msg := range msgChan {
		log.Println("Received message with routing key:", msg.RoutingKey)
		if handler, ok := dH.handlers[msg.RoutingKey]; ok {
			handler(&msg)

		} else {
			log.Printf("No handler found for routing key: %s", msg.RoutingKey)
		}
	}

}
func (dH *DecisionEventHandler) handleValuationCalculated(msg rabbitmq.Delivery) {
	var e event.ValuationCalculatedEvent
	err := json.Unmarshal(msg.Body, &e)
	if err != nil {
		log.Printf("Error unmarshalling ValuationCalculatedEvent: %v", err)
		return
	}
	fmt.Println("jestem w handle")
	_, err = dH.decisionService.PrepareDecision(e.ClaimID, e.PayoutAmount)
	if err != nil {
		log.Printf("err")
	}
}
