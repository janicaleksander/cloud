package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/notificationservice/domain"
	"github.com/mehdihadeli/go-mediatr"
)

type GetNotificationsForClaimIDQuery struct {
	ClaimID string
}

type GetNotificationsForClaimIDQueryResult struct {
	Notifications []*GetNotificationQueryResponse `json:"notifications"`
}
type GetNotificationsForClaimIDQueryHandler struct {
	repo domain.NotificationRepository
}

func NewGetNotificationsForClaimIDQueryHandler(repo domain.NotificationRepository) *GetNotificationsForClaimIDQueryHandler {
	return &GetNotificationsForClaimIDQueryHandler{repo: repo}
}

func (h *GetNotificationsForClaimIDQueryHandler) SelfRegister() error {
	return mediatr.RegisterRequestHandler[*GetNotificationsForClaimIDQuery, *GetNotificationsForClaimIDQueryResult](h)
}

func (h *GetNotificationsForClaimIDQueryHandler) Handle(ctx context.Context, query *GetNotificationsForClaimIDQuery) (*GetNotificationsForClaimIDQueryResult, error) {
	cid, err := uuid.Parse(query.ClaimID)
	if err != nil {
		return nil, err
	}
	notifications, err := h.repo.GetNotificationsByClaimID(ctx, cid)
	if err != nil {
		return nil, err
	}
	result := &GetNotificationsForClaimIDQueryResult{
		Notifications: make([]*GetNotificationQueryResponse, len(notifications)),
	}
	for i, n := range notifications {
		result.Notifications[i] = NotificationDomainToQueryResponse(n)
	}
	return result, nil
}
