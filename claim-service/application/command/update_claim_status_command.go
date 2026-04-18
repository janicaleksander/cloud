package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/claimservice/application/interfaces"
	"github.com/janicaleksander/cloud/claimservice/domain"
	"github.com/mehdihadeli/go-mediatr"
)

type UpdateClaimStatusCommand struct {
	ClaimID string
	Status  string
}

type UpdateClaimStatusCommandHandler struct {
	repo      domain.ClaimRepository
	publisher interfaces.ClaimEventPublisher
}

func NewUpdateClaimStatusCommandHandler(r domain.ClaimRepository, pub interfaces.ClaimEventPublisher) *UpdateClaimStatusCommandHandler {
	return &UpdateClaimStatusCommandHandler{
		repo:      r,
		publisher: pub,
	}
}

func (h *UpdateClaimStatusCommandHandler) SelfRegister() error {
	return mediatr.RegisterRequestHandler[*UpdateClaimStatusCommand, *mediatr.Unit](h)
}

func (h *UpdateClaimStatusCommandHandler) Handle(ctx context.Context, command *UpdateClaimStatusCommand) (*mediatr.Unit, error) {
	newStatus, err := domain.StringToStatus(command.Status)
	if err != nil {
		return nil, err
	}
	cid, err := uuid.Parse(command.ClaimID)
	if err != nil {
		return nil, err
	}
	err = h.repo.UpdateStatus(context.Background(), cid, newStatus)
	if err != nil {
		return nil, err
	}
	return &mediatr.Unit{}, nil
}
