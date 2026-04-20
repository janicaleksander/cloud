package command

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/common/event"
	"github.com/janicaleksander/cloud/decisionservice/application/interfaces"
	"github.com/janicaleksander/cloud/decisionservice/domain"
	"github.com/mehdihadeli/go-mediatr"
)

type UpdateDecisionStateCommand struct {
	DecisionID string
	NewState   string
	EmpID      string
	Reason     string
}

type UpdateDecisionStateCommandHandler struct {
	repo      domain.DecisionRepository
	publisher interfaces.DecisionPublisher
}

func NewUpdateDecisionStateCommandHandler(repo domain.DecisionRepository, p interfaces.DecisionPublisher) *UpdateDecisionStateCommandHandler {
	return &UpdateDecisionStateCommandHandler{repo: repo, publisher: p}
}

func (h *UpdateDecisionStateCommandHandler) SelfRegister() error {
	return mediatr.RegisterRequestHandler[*UpdateDecisionStateCommand, *mediatr.Unit](h)
}

func (h *UpdateDecisionStateCommandHandler) Handle(ctx context.Context, cmd *UpdateDecisionStateCommand) (*mediatr.Unit, error) {
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
	if oldDecisionDomain.Result != domain.WAITING {
		return nil, errors.New("already accepted/denied")
	}
	oldDecisionDomain.Result = domain.StringToResult(cmd.NewState)
	oldDecisionDomain.EmployeeID = eid
	h.makeDecision(oldDecisionDomain, cmd.Reason)
	return &mediatr.Unit{}, nil
}

func (h *UpdateDecisionStateCommandHandler) makeDecision(newDecision *domain.Decision, reason string) {
	slog.Info("Making decision for claim", "claimID", newDecision.ClaimID, "result", newDecision.Result, "employeeID", newDecision.EmployeeID)
	if newDecision.Result == domain.ACCEPTED {
		err := h.publisher.Publish("events", event.PayoutApprovedEvent{
			ClaimID:              newDecision.ClaimID.String(),
			AcceptedPayoutAmount: newDecision.Payout,
			ByEmployeeID:         newDecision.EmployeeID.String(),
		})
		if err != nil {
			slog.Error("Failed to publish PayoutApprovedEvent", "error", err)
			return
		}
		cmd := &UpdateEmpCommand{
			NewState:   string(newDecision.Result),
			DecisionID: newDecision.ID.String(),
			EmpID:      newDecision.EmployeeID.String(),
		}
		_, err = mediatr.Send[*UpdateEmpCommand, *mediatr.Unit](context.Background(), cmd)

		if err != nil {
			slog.Error("Failed to send UpdateDecisionStateCommand", "error", err)
			return
		}
	}
	if newDecision.Result == domain.REJECTED {
		err := h.publisher.Publish("events", event.PayoutRejectedEvent{
			ClaimID:      newDecision.ClaimID.String(),
			Reason:       reason,
			ByEmployeeID: newDecision.EmployeeID.String(),
		})
		if err != nil {
			slog.Error("Failed to publish PayoutRejectedEvent", "error", err)
			return
		}
		cmd := &UpdateEmpCommand{
			NewState:   string(newDecision.Result),
			DecisionID: newDecision.ID.String(),
			EmpID:      newDecision.EmployeeID.String(),
		}
		_, err = mediatr.Send[*UpdateEmpCommand, *mediatr.Unit](context.Background(), cmd)

		if err != nil {
			slog.Error("Failed to send UpdateDecisionStateCommand", "error", err)
			return
		}
	}

}
