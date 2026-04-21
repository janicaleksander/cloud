package command

import (
	"context"

	"github.com/janicaleksander/cloud/valuationservice/domain"
	"github.com/mehdihadeli/go-mediatr"
)

type CreateValuationCommand struct {
	ID      string
	ClaimID string
	Amount  float64
	Parts   []PartCommand
}

type PartCommand struct {
	ID   string
	Name string
	Cost float64
}
type CreateValuationCommandHandler struct {
	repo domain.ValuationRepository
}

func NewCreateValuationCommandHandler(repo domain.ValuationRepository) *CreateValuationCommandHandler {
	return &CreateValuationCommandHandler{repo: repo}
}

func (h *CreateValuationCommandHandler) SelfRegister() error {
	return mediatr.RegisterRequestHandler[*CreateValuationCommand, *mediatr.Unit](h)
}

func (h *CreateValuationCommandHandler) Handle(ctx context.Context, cmd *CreateValuationCommand) (*mediatr.Unit, error) {
	valuationDomain := CreateValuationCommandToDomain(cmd)
	_, err := h.repo.Save(ctx, valuationDomain)
	if err != nil {
		return nil, err
	}
	return &mediatr.Unit{}, nil

}
