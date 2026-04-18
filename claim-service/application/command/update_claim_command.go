package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/claimservice/application/interfaces"
	"github.com/janicaleksander/cloud/claimservice/domain"
	"github.com/janicaleksander/cloud/common/event"
	"github.com/mehdihadeli/go-mediatr"
)

type UpdateClaimCommand struct {
	ClaimID  string
	NewEmail string
}

type UpdateClaimCommandHandler struct {
	repo      domain.ClaimRepository
	publisher interfaces.ClaimEventPublisher
}

func NewUpdateClaimCommandHandler(r domain.ClaimRepository, pub interfaces.ClaimEventPublisher) *UpdateClaimCommandHandler {
	return &UpdateClaimCommandHandler{
		repo:      r,
		publisher: pub,
	}
}

func (h *UpdateClaimCommandHandler) SelfRegister() error {
	return mediatr.RegisterRequestHandler[*UpdateClaimCommand, *mediatr.Unit](h)
}

func (h *UpdateClaimCommandHandler) Handle(ctx context.Context, command *UpdateClaimCommand) (*mediatr.Unit, error) {
	if command.NewEmail == "" {
		return nil, nil
	}
	cid, err := uuid.Parse(command.ClaimID)
	if err != nil {
		return nil, err
	}
	oldClaimDomain, err := h.repo.GetById(context.Background(), cid)
	if err != nil {
		return nil, err
	}
	oldClaimDomain.Email = command.NewEmail
	_, err = h.repo.Update(context.Background(), oldClaimDomain)
	if err != nil {
		return nil, err
	}
	err = h.publisher.Publish("events", event.ChangeEmailForNotification{
		ClaimID: oldClaimDomain.ID.String(),
		Email:   command.NewEmail,
	})
	return &mediatr.Unit{}, err
}
