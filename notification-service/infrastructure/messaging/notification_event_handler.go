package messaging

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/janicaleksander/cloud/common/event"
	"github.com/janicaleksander/cloud/common/rabbitmq"
	"github.com/janicaleksander/cloud/common/rabbitmq/utils"
	"github.com/janicaleksander/cloud/notificationservice/application"
	"github.com/janicaleksander/cloud/notificationservice/domain"
)

const exchangeName = "events"
const queueName = "notification-service"

type NotificationEventHandler struct {
	notificationService *application.NotificationService
	emailService        *application.EmailService
	handlers            map[string]rabbitmq.HandlerFunc
}

func NewNotificationHandler(notificationService *application.NotificationService, emailService *application.EmailService) *NotificationEventHandler {
	slog.Info("Creating NotificationEventHandler")
	h := &NotificationEventHandler{
		notificationService: notificationService,
		emailService:        emailService,
		handlers:            make(map[string]rabbitmq.HandlerFunc),
	}
	h.registerHandlers()
	return h
}

func (n *NotificationEventHandler) Run(rabbit *rabbitmq.RabbitMQ) error {
	slog.Info("Running NotificationEventHandler")
	bindingKeys := make([]string, 0, len(n.handlers))
	for k := range n.handlers {
		bindingKeys = append(bindingKeys, k)
	}
	msgs, err := rabbitmq.SubscribeRaw(rabbit, exchangeName, queueName, bindingKeys...)
	if err != nil {
		return err
	}
	go n.dispatch(msgs)
	return nil
}

func (n *NotificationEventHandler) registerHandlers() {
	slog.Info("Registering event handlers for NotificationEventHandler")
	n.handlers[rabbitmq.RouteKeyToTopicNotation(utils.NameOfType(event.ClaimSubmittedEvent{}))] = n.handleClaimSubmitted
	n.handlers[rabbitmq.RouteKeyToTopicNotation(utils.NameOfType(event.PolicyVerifiedEvent{}))] = n.handlePolicyVerified
	n.handlers[rabbitmq.RouteKeyToTopicNotation(utils.NameOfType(event.PolicyDeniedEvent{}))] = n.handlePolicyDenied
	n.handlers[rabbitmq.RouteKeyToTopicNotation(utils.NameOfType(event.PayoutApprovedEvent{}))] = n.handlePayoutApproved
	n.handlers[rabbitmq.RouteKeyToTopicNotation(utils.NameOfType(event.PayoutRejectedEvent{}))] = n.handlePayoutRejected
	n.handlers[rabbitmq.RouteKeyToTopicNotation(utils.NameOfType(event.RegisterUserForNotificationEvent{}))] = n.handleRegisterUserForNotification
}

func (n *NotificationEventHandler) dispatch(msgs rabbitmq.MsgChan) {
	for msg := range msgs {
		if handler, ok := n.handlers[msg.RoutingKey]; ok {
			handler(&msg)
		} else {
			slog.Info("Received message with no handler", "routingKey", msg.RoutingKey)
		}
	}
}

func (n *NotificationEventHandler) sendAndSave(claimID uint, emailMsg application.EmailMessage) {
	slog.Info("Sending email and saving notification", "claimID", claimID, "to", emailMsg.To)
	if err := n.emailService.Send(emailMsg); err != nil {
		slog.Error("Error sending email", "claimID", claimID, "to", emailMsg.To, "error", err)
		return
	}
	_, err := n.notificationService.CreateNotification(&domain.Notification{
		ClaimID: claimID,
		Body:    emailMsg.Body,
		SentTo:  emailMsg.To,
		Time:    time.Now(),
	})
	if err != nil {
		slog.Error("Error saving notification", "claimID", claimID, "to", emailMsg.To, "error", err)
	}
}

func (n *NotificationEventHandler) handleClaimSubmitted(msg rabbitmq.Delivery) {
	var e event.ClaimSubmittedEvent
	if err := json.Unmarshal(msg.Body, &e); err != nil {
		slog.Error("Error unmarshalling ClaimSubmittedEvent", "error", err)
		return
	}
	slog.Info("Handling ClaimSubmittedEvent", "claimID", e.ClaimID, "userID", e.UserID, "vin", e.VIN)

	email, err := n.notificationService.GetEmailByClaimID(e.ClaimID)
	if err != nil {
		slog.Error("Error getting email for ClaimID", "claimID", e.ClaimID, "error", err)
		return
	}
	n.sendAndSave(e.ClaimID, application.EmailMessage{
		To:      email,
		Subject: fmt.Sprintf("Claim #%d Received", e.ClaimID),
		Body: fmt.Sprintf("Your insurance claim has been submitted successfully.\n\n"+
			"Claim ID: %d\nVehicle VIN: %s\nAccident Date: %s\nEvidence Files: %d\n\n"+
			"We will review your claim shortly and contact you with updates.",
			e.ClaimID, e.VIN, e.AccidentDate.Format("2006-01-02"), len(e.StorageURL)),
	})
}

func (n *NotificationEventHandler) handlePolicyVerified(msg rabbitmq.Delivery) {
	var e event.PolicyVerifiedEvent
	if err := json.Unmarshal(msg.Body, &e); err != nil {
		slog.Error("Error unmarshalling PolicyVerifiedEvent", "error", err)
		return
	}
	slog.Info("Handling PolicyVerifiedEvent", "claimID", e.ClaimID)

	email, err := n.notificationService.GetEmailByClaimID(e.ClaimID)
	if err != nil {
		slog.Error("Error getting email for ClaimID", "claimID", e.ClaimID, "error", err)
		return
	}
	n.sendAndSave(e.ClaimID, application.EmailMessage{
		To:      email,
		Subject: fmt.Sprintf("Claim #%d - Policy Verified", e.ClaimID),
		Body: fmt.Sprintf("Great news! We have verified your policy and confirmed coverage.\n\n"+
			"Claim ID: %d\nEvidence files reviewed: %d\n\n"+
			"Your claim is moving forward.",
			e.ClaimID, len(e.StorageURL)),
	})
}

func (n *NotificationEventHandler) handlePolicyDenied(msg rabbitmq.Delivery) {
	var e event.PolicyDeniedEvent
	if err := json.Unmarshal(msg.Body, &e); err != nil {
		slog.Error("Error unmarshalling PolicyDeniedEvent", "error", err)
		return
	}
	slog.Info("Handling PolicyDeniedEvent", "claimID", e.ClaimID, "reason", e.Reason)
	email, err := n.notificationService.GetEmailByClaimID(e.ClaimID)
	if err != nil {
		slog.Error("Error getting email for ClaimID", "claimID", e.ClaimID, "error", err)
		return
	}
	n.sendAndSave(e.ClaimID, application.EmailMessage{
		To:      email,
		Subject: fmt.Sprintf("Claim #%d - Policy Verification Failed", e.ClaimID),
		Body: fmt.Sprintf("Unfortunately, we were unable to verify your policy coverage.\n\n"+
			"Claim ID: %d\nReason: %s\n\n"+
			"Please contact our support team for more information.",
			e.ClaimID, e.Reason),
	})
}

func (n *NotificationEventHandler) handlePayoutApproved(msg rabbitmq.Delivery) {
	var e event.PayoutApprovedEvent
	if err := json.Unmarshal(msg.Body, &e); err != nil {
		slog.Error("Error unmarshalling PayoutApprovedEvent", "error", err)
		return
	}
	slog.Info("Handling PayoutApprovedEvent", "claimID", e.ClaimID, "approvedPayoutAmount", e.AcceptedPayoutAmount, "byEmployeeID", e.ByEmployeeID)
	email, err := n.notificationService.GetEmailByClaimID(e.ClaimID)
	if err != nil {
		slog.Error("Error getting email for ClaimID", "claimID", e.ClaimID, "error", err)
		return
	}
	n.sendAndSave(e.ClaimID, application.EmailMessage{
		To:      email,
		Subject: fmt.Sprintf("Claim #%d - Payout Approved", e.ClaimID),
		Body: fmt.Sprintf("Great news! Your claim has been approved.\n\n"+
			"Claim ID: %d\nApproved Payout Amount: $%.2f\nApproved by Employee ID: %d\n\n"+
			"The funds will be transferred within 3-5 business days.",
			e.ClaimID, e.AcceptedPayoutAmount, e.ByEmployeeID),
	})
}

func (n *NotificationEventHandler) handlePayoutRejected(msg rabbitmq.Delivery) {
	var e event.PayoutRejectedEvent
	if err := json.Unmarshal(msg.Body, &e); err != nil {
		slog.Error("Error unmarshalling PayoutRejectedEvent", "error", err)
		return
	}
	slog.Info("Handling PayoutRejectedEvent", "claimID", e.ClaimID, "reason", e.Reason, "byEmployeeID", e.ByEmployeeID)
	email, err := n.notificationService.GetEmailByClaimID(e.ClaimID)
	if err != nil {
		slog.Error("Error getting email for ClaimID", "claimID", e.ClaimID, "error", err)
		return
	}
	n.sendAndSave(e.ClaimID, application.EmailMessage{
		To:      email,
		Subject: fmt.Sprintf("Claim #%d - Payout Rejected", e.ClaimID),
		Body: fmt.Sprintf("We regret to inform you that your claim payout has been rejected.\n\n"+
			"Claim ID: %d\nReason: %s\nReviewed by Employee ID: %d\n\n"+
			"If you believe this decision is incorrect, please contact our support team.",
			e.ClaimID, e.Reason, e.ByEmployeeID),
	})
}

func (n *NotificationEventHandler) handleRegisterUserForNotification(msg rabbitmq.Delivery) {
	var e event.RegisterUserForNotificationEvent
	if err := json.Unmarshal(msg.Body, &e); err != nil {
		slog.Error("Error unmarshalling RegisterUserForNotificationEvent", "error", err)
		return
	}
	slog.Info("Handling RegisterUserForNotificationEvent", "claimID", e.ClaimID, "email", e.Email)
	_, err := n.notificationService.CreateNotificationReceiver(&domain.NotificationReceiver{
		ClaimID: e.ClaimID,
		Email:   e.Email,
	})
	if err != nil {
		slog.Error("Error creating notification receiver", "claimID", e.ClaimID, "email", e.Email, "error", err)
	}
}
