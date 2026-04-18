package command

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/policyverificationservice/domain"
	"github.com/mehdihadeli/go-mediatr"
)

type UpdatePolicyCommand struct {
	PolicyID string
	NewFrom  time.Time
	NewTo    time.Time
}

type UpdatePolicyCommandHandler struct {
	repo domain.PolicyRepository
}

func NewUpdatePolicyCommandHandler(repo domain.PolicyRepository) *UpdatePolicyCommandHandler {
	return &UpdatePolicyCommandHandler{repo: repo}
}
func (h *UpdatePolicyCommandHandler) SelfRegister() error {
	return mediatr.RegisterRequestHandler[*UpdatePolicyCommand, *mediatr.Unit](h)
}

func (h *UpdatePolicyCommandHandler) Handle(ctx context.Context, cmd *UpdatePolicyCommand) (*mediatr.Unit, error) {
	pid, err := uuid.Parse(cmd.PolicyID)
	if err != nil {
		return nil, err
	}
	oldPolicy, err := h.repo.GetById(ctx, pid)
	if err != nil {
		return nil, err
	}
	updated := *oldPolicy
	newFrom := cmd.NewFrom
	newTo := cmd.NewTo
	if newFrom != (time.Time{}) {
		updated.From = newFrom
	}
	if newTo != (time.Time{}) {
		updated.To = newTo
	}
	_, err = h.repo.Update(context.Background(), &updated)
	return &mediatr.Unit{}, err
}
