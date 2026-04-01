package persistance

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
