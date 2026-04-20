package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/decisionservice/domain"
	"github.com/mehdihadeli/go-mediatr"
)

type DeleteDecisionCommand struct {
	DecisionID string
}

type DeleteDecisionCommandHandler struct {
	repo domain.DecisionRepository
}

func NewDeleteDecisionCommandHandler(r domain.DecisionRepository) *DeleteDecisionCommandHandler {
	return &DeleteDecisionCommandHandler{repo: r}
}

func (h *DeleteDecisionCommandHandler) SelfRegister() error {
	return mediatr.RegisterRequestHandler[*DeleteDecisionCommand, *mediatr.Unit](h)
}

func (h *DeleteDecisionCommandHandler) Handle(ctx context.Context, cmd *DeleteDecisionCommand) (*mediatr.Unit, error) {
	did, err := uuid.Parse(cmd.DecisionID)
	if err != nil {
		return nil, err
	}
	err = h.repo.DeleteById(ctx, did)
	if err != nil {
		return nil, err
	}
	return &mediatr.Unit{}, nil
}
