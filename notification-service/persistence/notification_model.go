package persistence

import (
	"time"
)

type NotificationModel struct {
	ID      string    `dynamodbav:"notification_id"`
	ClaimID string    `dynamodbav:"claim_id"`
	Body    string    `dynamodbav:"body"`
	SentTo  string    `dynamodbav:"sent_to"`
	Time    time.Time `dynamodbav:"time"`
}

type NotificationReceiverModel struct {
	ID      string `dynamodbav:"notification_receiver_id"`
	ClaimID string `dynamodbav:"claim_id"`
	Email   string `dynamodbav:"email"`
}
