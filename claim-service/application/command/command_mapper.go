package command

import (
	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/claimservice/domain"
)

func CreateClaimCommandToDomain(cmd *CreateClaimCommand) *domain.Claim {
	cid, _ := uuid.Parse(cmd.ID)
	uid, _ := uuid.Parse(cmd.UserID)
	return &domain.Claim{
		ID:           cid,
		UserID:       uid,
		Email:        cmd.Email,
		VIN:          cmd.VIN,
		AccidentDate: cmd.AccidentDate,
	}
}
