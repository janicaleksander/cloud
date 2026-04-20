package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/notificationservice/domain"
	"github.com/mehdihadeli/go-mediatr"
)

type GetEmailByClaimIDQuery struct {
	ClaimID string
}

type GetEmailByClaimIDQueryResponse struct {
	Email string `json:"email"`
}
type GetEmailByClaimIDQueryHandler struct {
	repo domain.NotificationRepository
}

func NewGetEmailByClaimIDQueryHandler(notificationRepository domain.NotificationRepository) *GetEmailByClaimIDQueryHandler {
	return &GetEmailByClaimIDQueryHandler{
		repo: notificationRepository,
	}
}

func (h *GetEmailByClaimIDQueryHandler) SelfRegister() error {
	return mediatr.RegisterRequestHandler[*GetEmailByClaimIDQuery, *GetEmailByClaimIDQueryResponse](h)
}

func (h *GetEmailByClaimIDQueryHandler) Handle(ctx context.Context, query *GetEmailByClaimIDQuery) (*GetEmailByClaimIDQueryResponse, error) {
	cid, err := uuid.Parse(query.ClaimID)
	if err != nil {
		return nil, err
	}
	email, err := h.repo.GetEmailByClaimID(ctx, cid)
	if err != nil {
		return nil, err
	}
	return &GetEmailByClaimIDQueryResponse{Email: email}, nil

}
