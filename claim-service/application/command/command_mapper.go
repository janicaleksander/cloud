package command

import (
	"os"

	"github.com/janicaleksander/cloud/claimservice/domain"
)

func CreateClaimCommandToDomain(cmd *CreateClaimCommand) *domain.Claim {
	return &domain.Claim{
		UserID:       cmd.UserID,
		Email:        cmd.Email,
		VIN:          cmd.VIN,
		AccidentDate: cmd.AccidentDate,
		Files:        nil,
	}
}

func ClaimDomainToCommand(claim *domain.Claim, objectFiles []*os.File) *CreateClaimCommand {
	return &CreateClaimCommand{
		UserID:       claim.UserID,
		Email:        claim.Email,
		VIN:          claim.VIN,
		AccidentDate: claim.AccidentDate,
		ObjectFiles:  objectFiles,
	}
}
