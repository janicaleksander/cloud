package command

import (
	"github.com/janicaleksander/cloud/policyverificationservice/domain"
)

func CreatePolicyCommandToDomain(cmd *CreatePolicyCommand) *domain.Policy {
	return &domain.Policy{
		ID:     cmd.ID,
		UserID: cmd.UserID,
		VIN:    cmd.VIN,
		From:   cmd.From,
		To:     cmd.To,
	}
}
