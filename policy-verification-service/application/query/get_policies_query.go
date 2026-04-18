package query

import (
	"context"

	"github.com/janicaleksander/cloud/policyverificationservice/domain"
	"github.com/mehdihadeli/go-mediatr"
)

type GetPoliciesQuery struct{}

type GetPoliciesQueryResponse struct {
	Policies []*GetPolicyQueryResponse `json:"policies"`
}
type GetPoliciesQueryHandler struct {
	repo domain.PolicyRepository
}

func NewGetPoliciesQueryHandler(repo domain.PolicyRepository) *GetPoliciesQueryHandler {
	return &GetPoliciesQueryHandler{repo: repo}
}

func (h *GetPoliciesQueryHandler) SelfRegister() error {
	return mediatr.RegisterRequestHandler[*GetPoliciesQuery, *GetPoliciesQueryResponse](h)
}

func (h *GetPoliciesQueryHandler) Handle(ctx context.Context, query *GetPoliciesQuery) (*GetPoliciesQueryResponse, error) {
	policies, err := h.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	response := &GetPoliciesQueryResponse{
		Policies: make([]*GetPolicyQueryResponse, len(policies)),
	}

	for i, policy := range policies {
		response.Policies[i] = &GetPolicyQueryResponse{
			ID:     policy.ID.String(),
			UserID: policy.UserID.String(),
			VIN:    policy.VIN,
			From:   policy.From,
			To:     policy.To,
		}
	}

	return response, nil
}
