package persistence

import (
	"github.com/janicaleksander/cloud/notificationservice/domain"
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
		ID:      receiver.ID,
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
		ID:      notification.ID,
		ClaimID: notification.ClaimID,
		Body:    notification.Body,
		SentTo:  notification.SentTo,
		Time:    notification.Time,
	}
}
