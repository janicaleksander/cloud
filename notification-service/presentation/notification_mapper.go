package presentation

import "github.com/janicaleksander/cloud/notificationservice/domain"

func GetNotificationDomainToResponse(notification *domain.Notification) *GetNotificationResponseDTO {
	return &GetNotificationResponseDTO{
		ID:      notification.ID,
		ClaimID: notification.ClaimID,
		Body:    notification.Body,
		SentTo:  notification.SentTo,
		Time:    notification.Time,
	}
}
