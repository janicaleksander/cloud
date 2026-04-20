package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/valuationservice/domain"
	"github.com/mehdihadeli/go-mediatr"
)

type UpdateValuationCommand struct {
	ValuationID string
	NewAmount   float64
}

type UpdateValuationCommandHandler struct {
	repo domain.ValuationRepository
}

func NewUpdateValuationCommandHandler(repo domain.ValuationRepository) *UpdateValuationCommandHandler {
	return &UpdateValuationCommandHandler{repo: repo}
}

func (h *UpdateValuationCommandHandler) SelfRegister() error {
	return mediatr.RegisterRequestHandler[*UpdateValuationCommand, *mediatr.Unit](h)
}

// TODO remember about ID
// TODO do mappers in services
func (h *UpdateValuationCommandHandler) Handle(ctx context.Context, cmd *UpdateValuationCommand) (*mediatr.Unit, error) {
	vid, err := uuid.Parse(cmd.ValuationID)
	if err != nil {
		return nil, err
	}
	valuationDomain, err := h.repo.GetById(ctx, vid)
	if err != nil {
		return nil, err
	}
	if valuationDomain.Amount != cmd.NewAmount && cmd.NewAmount != 0 {
		_, err = h.repo.Update(ctx, &domain.Valuation{
			ID:      valuationDomain.ID,
			ClaimID: valuationDomain.ClaimID,
			Amount:  cmd.NewAmount,
			Parts:   valuationDomain.Parts,
		})
		if err != nil {
			return nil, err
		}
	}

	return &mediatr.Unit{}, nil
}
