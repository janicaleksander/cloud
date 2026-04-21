package query

import (
	"context"

	"github.com/janicaleksander/cloud/decisionservice/domain"
	"github.com/mehdihadeli/go-mediatr"
)

type GetWaitingDecisionsQuery struct{}

type GetWaitingDecisionsQueryResponse struct {
	Waiting []*GetDecisionQueryResult `json:"waiting"`
}
type GetWaitingDecisionsQueryHandler struct {
	repo domain.DecisionRepository
}

func NewGetWaitingDecisionsQueryHandler(r domain.DecisionRepository) *GetWaitingDecisionsQueryHandler {
	return &GetWaitingDecisionsQueryHandler{repo: r}
}

func (h *GetWaitingDecisionsQueryHandler) SelfRegister() error {
	return mediatr.RegisterRequestHandler[*GetWaitingDecisionsQuery, *GetWaitingDecisionsQueryResponse](h)
}

func (h *GetWaitingDecisionsQueryHandler) Handle(ctx context.Context, q *GetWaitingDecisionsQuery) (*GetWaitingDecisionsQueryResponse, error) {
	decisions, err := h.repo.GetAllWaiting(ctx)
	if err != nil {
		return nil, err
	}

	results := make([]*GetDecisionQueryResult, len(decisions))
	for i, d := range decisions {
		results[i] = DecisionDomainToQueryResponse(d)
	}
	return &GetWaitingDecisionsQueryResponse{Waiting: results}, nil
}
