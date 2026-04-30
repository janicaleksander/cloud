package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID      uuid.UUID
	ClaimID uuid.UUID
	Body    string
	SentTo  string
	Time    time.Time
}

type NotificationReceiver struct {
	ID      uuid.UUID
	ClaimID uuid.UUID
	Email   string
}

func NewNotification(id, claimID uuid.UUID, body, sentTo string, time time.Time) *Notification {
	return &Notification{
		ID:      id,
		ClaimID: claimID,
		Body:    body,
		SentTo:  sentTo,
		Time:    time,
	}
}

func NewNotificationReceiver(id, claimID uuid.UUID, email string) *NotificationReceiver {
	return &NotificationReceiver{
		ID:      id,
		ClaimID: claimID,
		Email:   email,
	}
}

type NotificationRepository interface {
	SaveNotification(context.Context, *Notification) (*Notification, error)
	GetNotification(context.Context, uuid.UUID) (*Notification, error)
	GetNotifications(context.Context) ([]*Notification, error)
	GetNotificationsByClaimID(context.Context, uuid.UUID) ([]*Notification, error)
	DeleteNotificationByID(context.Context, uuid.UUID) error

	SaveNotificationReceiver(context.Context, *NotificationReceiver) (*NotificationReceiver, error)
	GetEmailByClaimID(context.Context, uuid.UUID) (string, error)
}
