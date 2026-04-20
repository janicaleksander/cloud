package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/decisionservice/domain"
	"github.com/mehdihadeli/go-mediatr"
)

type UpdateEmpCommand struct {
	DecisionID string
	EmpID      string
	NewState   string
}

type UpdateEmpCommandHandler struct {
	repo domain.DecisionRepository
}

func NewUpdateEmpCommandHandler(repo domain.DecisionRepository) *UpdateEmpCommandHandler {
	return &UpdateEmpCommandHandler{repo: repo}
}

func (h *UpdateEmpCommandHandler) SelfRegister() error {
	return mediatr.RegisterRequestHandler[*UpdateEmpCommand, *mediatr.Unit](h)
}

func (h *UpdateEmpCommandHandler) Handle(ctx context.Context, cmd *UpdateEmpCommand) (*mediatr.Unit, error) {
	cid, err := uuid.Parse(cmd.DecisionID)
	if err != nil {
		return nil, err
	}
	oldDecisionDomain, err := h.repo.GetByID(ctx, cid)
	if err != nil {
		return nil, err
	}
	eid, err := uuid.Parse(cmd.EmpID)
	if err != nil {
		return nil, err
	}
	oldDecisionDomain.EmployeeID = eid
	oldDecisionDomain.Result = domain.StringToResult(cmd.NewState)
	_, err = h.repo.Update(ctx, oldDecisionDomain)
	if err != nil {
		return nil, err
	}
	return &mediatr.Unit{}, nil
}
