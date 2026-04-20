package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/valuationservice/domain"
	"github.com/mehdihadeli/go-mediatr"
)

type DeleteValuationCommand struct {
	ID string
}

type DeleteValuationCommandHandler struct {
	repo domain.ValuationRepository
}

func NewDeleteValuationCommandHandler(repo domain.ValuationRepository) *DeleteValuationCommandHandler {
	return &DeleteValuationCommandHandler{repo: repo}
}

func (h *DeleteValuationCommandHandler) SelfRegister() error {
	return mediatr.RegisterRequestHandler[*DeleteValuationCommand, *mediatr.Unit](h)
}

// TODO remember about ID
// TODO do mappers in services
func (h *DeleteValuationCommandHandler) Handle(ctx context.Context, cmd *DeleteValuationCommand) (*mediatr.Unit, error) {
	valuationID, err := uuid.Parse(cmd.ID)
	if err != nil {
		return nil, err
	}
	err = h.repo.DeleteById(ctx, valuationID)
	if err != nil {
		return nil, err
	}
	return &mediatr.Unit{}, nil
}
