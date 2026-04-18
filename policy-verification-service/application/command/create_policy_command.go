package command

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/policyverificationservice/domain"
	"github.com/mehdihadeli/go-mediatr"
)

type CreatePolicyCommand struct {
	ID     uuid.UUID
	UserID uuid.UUID
	VIN    string
	From   time.Time
	To     time.Time
}

type CreatePolicyCommandHandler struct {
	repo domain.PolicyRepository
}

func NewCreatePolicyCommandHandler(repo domain.PolicyRepository) *CreatePolicyCommandHandler {
	return &CreatePolicyCommandHandler{repo: repo}

}

func (h *CreatePolicyCommandHandler) SelfRegister() error {
	return mediatr.RegisterRequestHandler[*CreatePolicyCommand, *mediatr.Unit](h)
}

func (h *CreatePolicyCommandHandler) Handle(ctx context.Context, cmd *CreatePolicyCommand) (*mediatr.Unit, error) {
	policyDomain := CreatePolicyCommandToDomain(cmd)
	_, err := h.repo.Save(ctx, policyDomain)
	if err != nil {
		return nil, err
	}
	return &mediatr.Unit{}, nil

}
