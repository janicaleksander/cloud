package command

import (
	"context"

	"github.com/janicaleksander/cloud/claimservice/domain"
	"github.com/mehdihadeli/go-mediatr"
)

type DeleteClaimCommand struct {
	ClaimID uint
}

type DeleteClaimCommandHandler struct {
	repo domain.ClaimRepository
}

func NewDeleteClaimCommandHandler(r domain.ClaimRepository) *DeleteClaimCommandHandler {
	return &DeleteClaimCommandHandler{repo: r}
}

func (h *DeleteClaimCommandHandler) SelfRegister() error {
	return mediatr.RegisterRequestHandler[*DeleteClaimCommand, *mediatr.Unit](h)
}

func (h *DeleteClaimCommandHandler) Handle(ctx context.Context, command *DeleteClaimCommand) (*mediatr.Unit, error) {
	err := h.repo.DeleteById(context.Background(), command.ClaimID)
	return &mediatr.Unit{}, err
}
