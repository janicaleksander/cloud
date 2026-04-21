package command

import (
	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/decisionservice/domain"
)

func PrepareDecisionCommandToDomain(cmd *PrepareDecisionCommand) *domain.Decision {
	cid, err := uuid.Parse(cmd.ClaimID)
	if err != nil {
		return nil
	}
	did, err := uuid.Parse(cmd.ID)
	if err != nil {
		return nil
	}
	return &domain.Decision{
		ID:      did,
		ClaimID: cid,
		Payout:  cmd.PayoutAmount,
		Result:  domain.WAITING,
	}
}
