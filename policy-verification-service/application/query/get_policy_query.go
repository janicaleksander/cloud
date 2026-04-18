package query

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/policyverificationservice/domain"
	"github.com/mehdihadeli/go-mediatr"
)

type GetPolicyQuery struct {
	PolicyID string
}

type GetPolicyQueryHandler struct {
	repo domain.PolicyRepository
}

type GetPolicyQueryResponse struct {
	ID     string    `json:"id"`
	UserID string    `json:"user_id"`
	VIN    string    `json:"vin"`
	From   time.Time `json:"from"`
	To     time.Time `json:"to"`
}

func NewGetPolicyQueryHandler(repo domain.PolicyRepository) *GetPolicyQueryHandler {
	return &GetPolicyQueryHandler{repo: repo}
}

func (h *GetPolicyQueryHandler) SelfRegister() error {
	return mediatr.RegisterRequestHandler[*GetPolicyQuery, *GetPolicyQueryResponse](h)
}

func (h *GetPolicyQueryHandler) Handle(ctx context.Context, query *GetPolicyQuery) (*GetPolicyQueryResponse, error) {
	pid, err := uuid.Parse(query.PolicyID)
	if err != nil {
		return nil, err
	}
	policyDomain, err := h.repo.GetById(ctx, pid)
	if err != nil {
		return nil, err
	}

	return &GetPolicyQueryResponse{
		ID:     policyDomain.ID.String(),
		UserID: policyDomain.UserID.String(),
		VIN:    policyDomain.VIN,
		From:   policyDomain.From,
		To:     policyDomain.To,
	}, nil

}
