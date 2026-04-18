package command

import (
	"os"

	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/claimservice/domain"
)

func CreateClaimCommandToDomain(cmd *CreateClaimCommand) *domain.Claim {
	uid, _ := uuid.Parse(cmd.UserID)
	return &domain.Claim{
		UserID:       uid,
		Email:        cmd.Email,
		VIN:          cmd.VIN,
		AccidentDate: cmd.AccidentDate,
	}
}

func ClaimDomainToCreateClaimCommand(claim *domain.Claim, objectFiles []*os.File) *CreateClaimCommand {
	return &CreateClaimCommand{
		UserID:       claim.UserID.String(),
		Email:        claim.Email,
		VIN:          claim.VIN,
		AccidentDate: claim.AccidentDate,
		ObjectFiles:  objectFiles,
	}
}
