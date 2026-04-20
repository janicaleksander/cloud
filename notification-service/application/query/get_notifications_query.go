package query

import (
	"context"

	"github.com/janicaleksander/cloud/notificationservice/domain"
	"github.com/mehdihadeli/go-mediatr"
)

type GetNotificationsQuery struct{}

type GetNotificationsQueryResponse struct {
	Notifications []*GetNotificationQueryResponse `json:"notifications"`
}

type GetNotificationsQueryHandler struct {
	repo domain.NotificationRepository
}

func NewGetNotificationsQueryHandler(repo domain.NotificationRepository) *GetNotificationsQueryHandler {
	return &GetNotificationsQueryHandler{repo: repo}
}

func (h *GetNotificationsQueryHandler) SelfRegister() error {
	return mediatr.RegisterRequestHandler[*GetNotificationsQuery, *GetNotificationsQueryResponse](h)
}

func (h *GetNotificationsQueryHandler) Handle(ctx context.Context, query *GetNotificationsQuery) (*GetNotificationsQueryResponse, error) {
	notifications, err := h.repo.GetNotifications(ctx)
	if err != nil {
		return nil, err
	}

	response := &GetNotificationsQueryResponse{
		Notifications: make([]*GetNotificationQueryResponse, len(notifications)),
	}

	for i, n := range notifications {
		response.Notifications[i] = &GetNotificationQueryResponse{
			ID:      n.ID.String(),
			ClaimID: n.ClaimID.String(),
			Body:    n.Body,
			SentTo:  n.SentTo,
			Time:    n.Time,
		}
	}

	return response, nil
}
