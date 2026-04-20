package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/decisionservice/domain"
	"github.com/mehdihadeli/go-mediatr"
)

type PrepareDecisionCommand struct {
	ID           string
	ClaimID      string
	PayoutAmount float64
}

type PrepareDecisionCommandHandler struct {
	repo domain.DecisionRepository
}

func NewPrepareDecisionCommandHandler(repo domain.DecisionRepository) *PrepareDecisionCommandHandler {
	return &PrepareDecisionCommandHandler{repo: repo}
}

func (h *PrepareDecisionCommandHandler) SelfRegister() error {
	return mediatr.RegisterRequestHandler[*PrepareDecisionCommand, *mediatr.Unit](h)
}

func (h *PrepareDecisionCommandHandler) Handle(ctx context.Context, cmd *PrepareDecisionCommand) (*mediatr.Unit, error) {
	cid, err := uuid.Parse(cmd.ClaimID)
	if err != nil {
		return nil, err
	}
	did, err := uuid.Parse(cmd.ID)
	if err != nil {
		return nil, err
	}
	decisionDomain := &domain.Decision{
		ID:      did,
		ClaimID: cid,
		Payout:  cmd.PayoutAmount,
		Result:  domain.WAITING,
	}
	_, err = h.repo.Save(ctx, decisionDomain)
	if err != nil {
		return nil, err
	}
	return &mediatr.Unit{}, nil
}
