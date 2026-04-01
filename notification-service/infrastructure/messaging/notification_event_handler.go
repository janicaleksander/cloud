package messaging

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/janicaleksander/cloud/common/event"
	"github.com/janicaleksander/cloud/common/rabbitmq"
	"github.com/janicaleksander/cloud/common/rabbitmq/utils"
	"github.com/janicaleksander/cloud/notificationservice/application"
	"github.com/janicaleksander/cloud/notificationservice/domain"
)

type NotificationEventHandler struct {
	notificationService *application.NotificationService
	emailService        *application.EmailService
	handlers            map[string]rabbitmq.HandlerFunc
}

func NewNotificationHandler(notificationService *application.NotificationService, emailService *application.EmailService) *NotificationEventHandler {
	h := &NotificationEventHandler{
		notificationService: notificationService,
		emailService:        emailService,
		handlers:            make(map[string]rabbitmq.HandlerFunc),
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
	n.handlers[rabbitmq.RouteKeyToTopicNotation(utils.NameOfType(event.RegisterUserForNotificationEvent{}))] = n.handleRegisterUserForNotification
	n.handlers[rabbitmq.RouteKeyToTopicNotation(utils.NameOfType(event.ChangeEmailForNotification{}))] = n.handleChangeEmailForNotification

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

	email, err := n.notificationService.GetEmailByClaimID(e.ClaimID)
	if err != nil {
		log.Printf("Error getting email for ClaimID %d: %v\n", e.ClaimID, err)
		return
	}
	emailMsg := application.EmailMessage{
		To:      email,
		Subject: fmt.Sprintf("Claim #%d Received", e.ClaimID),
		Body: fmt.Sprintf("Your insurance claim has been submitted successfully.\n\n"+
			"Claim ID: %d\nVehicle VIN: %s\nAccident Date: %s\nEvidence Files: %d\n\n"+
			"We will review your claim shortly and contact you with updates.",
			e.ClaimID, e.VIN, e.AccidentDate.Format("2006-01-02"), len(e.StorageURL)),
	}
	if err := n.emailService.Send(emailMsg); err != nil {
		log.Printf("Error sending email for ClaimSubmittedEvent: %v\n", err)
	}
}

func (n *NotificationEventHandler) handlePolicyVerified(msg rabbitmq.Delivery) {
	var e event.PolicyVerifiedEvent
	err := json.Unmarshal(msg.Body, &e)
	if err != nil {
		log.Println("Error unmarshalling PolicyVerifiedEvent:", err)
		return
	}
	log.Printf("[NOTIFICATION] Policy Verified - ClaimID: %d\n", e.ClaimID)

	email, err := n.notificationService.GetEmailByClaimID(e.ClaimID)
	if err != nil {
		log.Printf("Error getting email for ClaimID %d: %v\n", e.ClaimID, err)
		return
	}
	emailMsg := application.EmailMessage{
		To:      email,
		Subject: fmt.Sprintf("Claim #%d - Policy Verified", e.ClaimID),
		Body: fmt.Sprintf("Great news! We have verified your policy and confirmed coverage.\n\n"+
			"Claim ID: %d\nEvidence files reviewed: %d\n\n"+
			"Your claim is moving forward.",
			e.ClaimID, len(e.StorageURL)),
	}
	if err := n.emailService.Send(emailMsg); err != nil {
		log.Printf("Error sending email for PolicyVerifiedEvent: %v\n", err)
	}
}

func (n *NotificationEventHandler) handlePolicyDenied(msg rabbitmq.Delivery) {
	var e event.PolicyDeniedEvent
	err := json.Unmarshal(msg.Body, &e)
	if err != nil {
		log.Println("Error unmarshalling PolicyDeniedEvent:", err)
		return
	}
	log.Printf("[NOTIFICATION] Policy Denied - ClaimID: %d, Reason: %s\n", e.ClaimID, e.Reason)

	email, err := n.notificationService.GetEmailByClaimID(e.ClaimID)
	if err != nil {
		log.Printf("Error getting email for ClaimID %d: %v\n", e.ClaimID, err)
		return
	}
	emailMsg := application.EmailMessage{
		To:      email,
		Subject: fmt.Sprintf("Claim #%d - Policy Verification Failed", e.ClaimID),
		Body: fmt.Sprintf("Unfortunately, we were unable to verify your policy coverage.\n\n"+
			"Claim ID: %d\nReason: %s\n\n"+
			"Please contact our support team for more information.",
			e.ClaimID, e.Reason),
	}
	if err := n.emailService.Send(emailMsg); err != nil {
		log.Printf("Error sending email for PolicyDeniedEvent: %v\n", err)
	}
}

func (n *NotificationEventHandler) handlePayoutApproved(msg rabbitmq.Delivery) {
	var e event.PayoutApprovedEvent
	err := json.Unmarshal(msg.Body, &e)
	if err != nil {
		log.Println("Error unmarshalling PayoutApprovedEvent:", err)
		return
	}
	log.Printf("[NOTIFICATION] Payout Approved - ClaimID: %d, Amount: %.2f\n", e.ClaimID, e.AcceptedPayoutAmount)

	email, err := n.notificationService.GetEmailByClaimID(e.ClaimID)
	if err != nil {
		log.Printf("Error getting email for ClaimID %d: %v\n", e.ClaimID, err)
		return
	}
	emailMsg := application.EmailMessage{
		To:      email,
		Subject: fmt.Sprintf("Claim #%d - Payout Approved", e.ClaimID),
		Body: fmt.Sprintf("Great news! Your claim has been approved.\n\n"+
			"Claim ID: %d\nApproved Payout Amount: $%.2f\nApproved by Employee ID: %d\n\n"+
			"The funds will be transferred within 3-5 business days.",
			e.ClaimID, e.AcceptedPayoutAmount, e.ByEmployeeID),
	}
	if err := n.emailService.Send(emailMsg); err != nil {
		log.Printf("Error sending email for PayoutApprovedEvent: %v\n", err)
	}
}

func (n *NotificationEventHandler) handlePayoutRejected(msg rabbitmq.Delivery) {
	var e event.PayoutRejectedEvent
	err := json.Unmarshal(msg.Body, &e)
	if err != nil {
		log.Println("Error unmarshalling PayoutRejectedEvent:", err)
		return
	}
	log.Printf("[NOTIFICATION] Payout Rejected - ClaimID: %d, Reason: %s\n", e.ClaimID, e.Reason)

	email, err := n.notificationService.GetEmailByClaimID(e.ClaimID)
	if err != nil {
		log.Printf("Error getting email for ClaimID %d: %v\n", e.ClaimID, err)
		return
	}
	emailMsg := application.EmailMessage{
		To:      email,
		Subject: fmt.Sprintf("Claim #%d - Payout Rejected", e.ClaimID),
		Body: fmt.Sprintf("We regret to inform you that your claim payout has been rejected.\n\n"+
			"Claim ID: %d\nReason: %s\nReviewed by Employee ID: %d\n\n"+
			"If you believe this decision is incorrect, please contact our support team.",
			e.ClaimID, e.Reason, e.ByEmployeeID),
	}
	if err := n.emailService.Send(emailMsg); err != nil {
		log.Printf("Error sending email for PayoutRejectedEvent: %v\n", err)
	}
}

func (n *NotificationEventHandler) handleRegisterUserForNotification(msg rabbitmq.Delivery) {
	//TODO: Implement handler for RegisterUserForNotificationEvent
	var e event.RegisterUserForNotificationEvent
	err := json.Unmarshal(msg.Body, &e)
	if err != nil {
		log.Println("Error unmarshalling RegisterUserForNotificationEvent:", err)
		return
	}
	log.Printf("[NOTIFICATION] Register User for Notification - ClaimID: %d, Email: %s\n", e.ClaimID, e.Email)
	err = n.notificationService.CreateNotificationReceiver(&domain.NotificationReceiver{
		ClaimID: e.ClaimID,
		Email:   e.Email,
	})
	if err != nil {
		log.Printf("Error creating notification receiver: %v\n", err)
	}

	// Here you would typically save the email and claim ID to your database
	// so that you can send notifications to this email when events related to this claim occur.
}

func (n *NotificationEventHandler) handleChangeEmailForNotification(msg rabbitmq.Delivery) {
	var e event.ChangeEmailForNotification
	err := json.Unmarshal(msg.Body, &e)
	if err != nil {
		log.Println("Error unmarshalling ChangeEmailForNotification:", err)
		return
	}
	log.Printf("[NOTIFICATION] Change Email for Notification - ClaimID: %d, New Email: %s\n", e.ClaimID, e.Email)
	err = n.notificationService.UpdateNotificationReceiver(&domain.NotificationReceiver{
		ClaimID: e.ClaimID,
		Email:   e.Email,
	})
	if err != nil {
		log.Printf("Error updating notification receiver: %v\n", err)
	}
	// Here you would typically update the email associated with the claim ID in your database

}
