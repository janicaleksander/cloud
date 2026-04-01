package messaging

import (
	"encoding/json"
	"log"

	"github.com/janicaleksander/cloud/common/event"
	"github.com/janicaleksander/cloud/common/rabbitmq"
	"github.com/janicaleksander/cloud/common/rabbitmq/utils"
)

type NotificationEventHandler struct {
	//todo notification service
	handlers map[string]rabbitmq.HandlerFunc
}

func NewNotificationHandler() *NotificationEventHandler {
	h := &NotificationEventHandler{
		handlers: make(map[string]rabbitmq.HandlerFunc),
	}
	h.registerHandlers()
	return h
}
func (n *NotificationEventHandler) Run(rabbit *rabbitmq.RabbitMQ) error {
	bindingKeys := make([]string, 0, len(n.handlers))
	for k := range n.handlers {
		bindingKeys = append(bindingKeys, k)
	}
	msgs, err := rabbitmq.SubscribeRaw(rabbit, "events", "notification-service", bindingKeys...)
	if err != nil {
		return err
	}
	go n.dispatch(msgs)
	return nil
}

func (n *NotificationEventHandler) registerHandlers() {
	n.handlers[rabbitmq.RouteKeyToTopicNotation(utils.NameOfType(event.ClaimSubmittedEvent{}))] = n.handleClaimSubmitted
	n.handlers[rabbitmq.RouteKeyToTopicNotation(utils.NameOfType(event.PolicyVerifiedEvent{}))] = n.handlePolicyVerified
	n.handlers[rabbitmq.RouteKeyToTopicNotation(utils.NameOfType(event.PolicyDeniedEvent{}))] = n.handlePolicyDenied
	n.handlers[rabbitmq.RouteKeyToTopicNotation(utils.NameOfType(event.PayoutApprovedEvent{}))] = n.handlePayoutApproved
	n.handlers[rabbitmq.RouteKeyToTopicNotation(utils.NameOfType(event.PayoutRejectedEvent{}))] = n.handlePayoutRejected

}

func (n *NotificationEventHandler) dispatch(msgs rabbitmq.MsgChan) {
	for msg := range msgs {
		if handler, ok := n.handlers[msg.RoutingKey]; ok {
			handler(&msg)
		} else {
			log.Println("no handler for routing key: ", msg.RoutingKey)
		}
	}
}

func (n *NotificationEventHandler) handleClaimSubmitted(msg rabbitmq.Delivery) {
	var e event.ClaimSubmittedEvent
	err := json.Unmarshal(msg.Body, &e)
	if err != nil {
		log.Println("Error unmarshalling ClaimSubmittedEvent:", err)
		return
	}
	log.Printf("[NOTIFICATION] Claim Submitted - ClaimID: %d, UserID: %d, VIN: %s, AccidentDate: %s, Evidence: %d files\n",
		e.ClaimID, e.UserID, e.VIN, e.AccidentDate.Format("2006-01-02"), len(e.StorageURL))
}

func (n *NotificationEventHandler) handlePolicyVerified(msg rabbitmq.Delivery) {
	var e event.PolicyVerifiedEvent
	err := json.Unmarshal(msg.Body, &e)
	if err != nil {
		log.Println("Error unmarshalling PolicyVerifiedEvent:", err)
		return
	}
	log.Printf("[NOTIFICATION] Policy Verified - ClaimID: %d, Evidence URLs: %d\n", e.ClaimID, len(e.StorageURL))
}

func (n *NotificationEventHandler) handlePolicyDenied(msg rabbitmq.Delivery) {
	var e event.PolicyDeniedEvent
	err := json.Unmarshal(msg.Body, &e)
	if err != nil {
		log.Println("Error unmarshalling PolicyDeniedEvent:", err)
		return
	}
	log.Printf("[NOTIFICATION] Policy Denied - ClaimID: %d, Reason: %s\n", e.ClaimID, e.Reason)
}

func (n *NotificationEventHandler) handlePayoutApproved(msg rabbitmq.Delivery) {
	var e event.PayoutApprovedEvent
	err := json.Unmarshal(msg.Body, &e)
	if err != nil {
		log.Println("Error unmarshalling PayoutApprovedEvent:", err)
		return
	}
	log.Printf("[NOTIFICATION] Payout Approved - ClaimID: %d, Amount: $%.2f\n", e.ClaimID, e.AcceptedPayoutAmount)
}

func (n *NotificationEventHandler) handlePayoutRejected(msg rabbitmq.Delivery) {
	var e event.PayoutRejectedEvent
	err := json.Unmarshal(msg.Body, &e)
	if err != nil {
		log.Println("Error unmarshalling PayoutRejectedEvent:", err)
		return
	}
	log.Printf("[NOTIFICATION] Payout Rejected - ClaimID: %d, Reason: %s\n", e.ClaimID, e.Reason)
}
