package query

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/notificationservice/domain"
	"github.com/mehdihadeli/go-mediatr"
)

type GetNotificationQuery struct {
	NotificationID string
}

type GetNotificationQueryResponse struct {
	ID      string    `json:"id"`
	ClaimID string    `json:"claim_id"`
	Body    string    `json:"body"`
	SentTo  string    `json:"sent_to"`
	Time    time.Time `json:"time"`
}

type GetNotificationQueryHandler struct {
	repo domain.NotificationRepository
}

func NewGetNotificationQueryHandler(repo domain.NotificationRepository) *GetNotificationQueryHandler {
	return &GetNotificationQueryHandler{repo: repo}
}

func (h *GetNotificationQueryHandler) SelfRegister() error {
	return mediatr.RegisterRequestHandler[*GetNotificationQuery, *GetNotificationQueryResponse](h)
}

func (h *GetNotificationQueryHandler) Handle(ctx context.Context, query *GetNotificationQuery) (*GetNotificationQueryResponse, error) {
	nid, err := uuid.Parse(query.NotificationID)
	if err != nil {
		return nil, err
	}
	notificationDomain, err := h.repo.GetNotification(ctx, nid)
	if err != nil {
		return nil, err
	}
	return NotificationDomainToQueryResponse(notificationDomain), nil
}
