package query

import (
	"context"

	"github.com/janicaleksander/cloud/decisionservice/domain"
	"github.com/mehdihadeli/go-mediatr"
)

type GetWaitingDecisionsQuery struct{}

type GetWaitingDecisionsQueryResponse struct {
	Waiting []*GetDecisionQueryResult
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
		results[i] = &GetDecisionQueryResult{
			ID:         d.ID.String(),
			EmployeeID: d.EmployeeID.String(),
			ClaimID:    d.ClaimID.String(),
			Payout:     d.Payout,
			State:      string(d.Result),
		}
	}
	return &GetWaitingDecisionsQueryResponse{Waiting: results}, nil
}
