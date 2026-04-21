package query

import "github.com/janicaleksander/cloud/notificationservice/domain"

func NotificationDomainToQueryResponse(notificationDomain *domain.Notification) *GetNotificationQueryResponse {
	return &GetNotificationQueryResponse{
		ID:      notificationDomain.ID.String(),
		ClaimID: notificationDomain.ClaimID.String(),
		Body:    notificationDomain.Body,
		SentTo:  notificationDomain.SentTo,
		Time:    notificationDomain.Time,
	}
}
