package query

import (
	"context"

	"github.com/janicaleksander/cloud/decisionservice/domain"
	"github.com/mehdihadeli/go-mediatr"
)

type GetDecisionsQuery struct{}

type GetDecisionsQueryResult struct {
	Decisions []*GetDecisionQueryResult `json:"decisions"`
}
type GetDecisionsQueryHandler struct {
	repo domain.DecisionRepository
}

func NewGetDecisionsQueryHandler(repo domain.DecisionRepository) *GetDecisionsQueryHandler {
	return &GetDecisionsQueryHandler{repo: repo}
}

func (h *GetDecisionsQueryHandler) SelfRegister() error {
	return mediatr.RegisterRequestHandler[*GetDecisionsQuery, *GetDecisionsQueryResult](h)
}

func (h *GetDecisionsQueryHandler) Handle(ctx context.Context, query *GetDecisionsQuery) (*GetDecisionsQueryResult, error) {
	decisions, err := h.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	result := &GetDecisionsQueryResult{
		Decisions: make([]*GetDecisionQueryResult, len(decisions)),
	}
	for i, d := range decisions {
		result.Decisions[i] = DecisionDomainToQueryResponse(d)
	}
	return result, nil

}
