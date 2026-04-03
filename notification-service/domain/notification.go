package domain

import (
	"context"
	"time"
)

type Notification struct {
	ID      uint
	ClaimID uint
	Body    string
	SentTo  string
	Time    time.Time
}

type NotificationReceiver struct {
	ID      uint
	ClaimID uint
	Email   string
}

type NotificationRepository interface {
	SaveNotification(context.Context, *Notification) (*Notification, error)
	GetNotification(context.Context, uint) (*Notification, error)
	GetNotifications(context.Context) ([]*Notification, error)
	GetNotificationsByClaimID(context.Context, uint) ([]*Notification, error)
	DeleteNotificationByID(context.Context, uint) error

	SaveNotificationReceiver(context.Context, *NotificationReceiver) (*NotificationReceiver, error)
	UpdateNotificationReceiver(context.Context, *NotificationReceiver) (*NotificationReceiver, error)
	GetEmailByClaimID(context.Context, uint) (string, error)
}
