package persistence

import (
	"github.com/janicaleksander/cloud/notificationservice/domain"
	"gorm.io/gorm"
)

func NotificationReceiverModelToDomain(receiver *NotificationReceiverModel) *domain.NotificationReceiver {
	return &domain.NotificationReceiver{
		ID:      receiver.ID,
		ClaimID: receiver.ClaimID,
		Email:   receiver.Email,
	}
}

func NotificationReceiverDomainToModel(receiver *domain.NotificationReceiver) *NotificationReceiverModel {
	return &NotificationReceiverModel{
		Model:   gorm.Model{ID: receiver.ID},
		ClaimID: receiver.ClaimID,
		Email:   receiver.Email,
	}
}

func NotificationModelToDomain(notification *NotificationModel) *domain.Notification {
	return &domain.Notification{
		ID:      notification.ID,
		ClaimID: notification.ClaimID,
		Body:    notification.Body,
		SentTo:  notification.SentTo,
		Time:    notification.Time,
	}
}

func NotificationDomainToModel(notification *domain.Notification) *NotificationModel {
	return &NotificationModel{
		Model:   gorm.Model{ID: notification.ID},
		ClaimID: notification.ClaimID,
		Body:    notification.Body,
		SentTo:  notification.SentTo,
		Time:    notification.Time,
	}
}
