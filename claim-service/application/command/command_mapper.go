package command

import (
	"os"

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

func ClaimDomainToCreateClaimCommand(claim *domain.Claim, objectFiles []*os.File) *CreateClaimCommand {
	return &CreateClaimCommand{
		ID:           claim.ID.String(),
		UserID:       claim.UserID.String(),
		Email:        claim.Email,
		VIN:          claim.VIN,
		AccidentDate: claim.AccidentDate,
		ObjectFiles:  objectFiles,
	}
}
