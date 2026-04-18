package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/policyverificationservice/domain"
	"github.com/mehdihadeli/go-mediatr"
)

type DeletePolicyCommand struct {
	PolicyID string
}

type DeletePolicyCommandHandler struct {
	repo domain.PolicyRepository
}

func NewDeletePolicyCommandHandler(repo domain.PolicyRepository) *DeletePolicyCommandHandler {
	return &DeletePolicyCommandHandler{repo: repo}
}

func (h *DeletePolicyCommandHandler) SelfRegister() error {
	return mediatr.RegisterRequestHandler[*DeletePolicyCommand, *mediatr.Unit](h)
}

func (h *DeletePolicyCommandHandler) Handle(ctx context.Context, cmd *DeletePolicyCommand) (*mediatr.Unit, error) {
	pid, err := uuid.Parse(cmd.PolicyID)
	if err != nil {
		return nil, err
	}
	err = h.repo.DeleteById(ctx, pid)
	if err != nil {
		return nil, err
	}
	return &mediatr.Unit{}, nil
}
