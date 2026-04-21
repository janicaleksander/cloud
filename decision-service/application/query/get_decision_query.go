package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/decisionservice/domain"
	"github.com/mehdihadeli/go-mediatr"
)

type GetDecisionQuery struct {
	DecisionID string
}

type GetDecisionQueryResult struct {
	ID         string  `json:"id"`
	EmployeeID string  `json:"employee_id,omitempty"`
	ClaimID    string  `json:"claim_id"`
	Payout     float64 `json:"payout"`
	State      string  `json:"state"`
}

type GetDecisionQueryHandler struct {
	repo domain.DecisionRepository
}

func NewGetDecisionQueryHandler(repo domain.DecisionRepository) *GetDecisionQueryHandler {
	return &GetDecisionQueryHandler{repo: repo}
}

func (h *GetDecisionQueryHandler) SelfRegister() error {
	return mediatr.RegisterRequestHandler[*GetDecisionQuery, *GetDecisionQueryResult](h)

}

func (h *GetDecisionQueryHandler) Handle(ctx context.Context, query *GetDecisionQuery) (*GetDecisionQueryResult, error) {
	did, err := uuid.Parse(query.DecisionID)
	if err != nil {
		return nil, err
	}
	decision, err := h.repo.GetByID(ctx, did)
	if err != nil {
		return nil, err
	}
	return DecisionDomainToQueryResponse(decision), nil
}
