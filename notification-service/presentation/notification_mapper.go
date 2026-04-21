package presentation

import (
	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/notificationservice/application/command"
	"github.com/janicaleksander/cloud/notificationservice/application/query"
)

func GetNotificationHTTPToQuery(notificationID uuid.UUID) *query.GetNotificationQuery {
	return &query.GetNotificationQuery{
		NotificationID: notificationID.String(),
	}
}

func GetNotificationsHTTPToQuery() *query.GetNotificationsQuery {
	return &query.GetNotificationsQuery{}
}

func GetNotificationsForClaimIDHTTPToQuery(claimID uuid.UUID) *query.GetNotificationsForClaimIDQuery {
	return &query.GetNotificationsForClaimIDQuery{
		ClaimID: claimID.String(),
	}
}

func DeleteNotificationHTTPToCommand(notificationID uuid.UUID) *command.DeleteNotificationCommand {
	return &command.DeleteNotificationCommand{
		NotificationID: notificationID.String(),
	}
}
