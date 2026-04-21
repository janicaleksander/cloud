package command

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/policyverificationservice/domain"
	"github.com/mehdihadeli/go-mediatr"
)

type CreatePolicyCommand struct {
	ID     string
	UserID string
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
	policyID, err := uuid.Parse(cmd.ID)
	if err != nil {
		return nil, err
	}
	userID, err := uuid.Parse(cmd.UserID)
	if err != nil {
		return nil, err
	}
	policyDomain := domain.NewPolicy(
		policyID,
		userID,
		cmd.VIN,
		cmd.From,
		cmd.To,
	)
	_, err = h.repo.Save(ctx, policyDomain)
	if err != nil {
		return nil, err
	}
	return &mediatr.Unit{}, nil

}
