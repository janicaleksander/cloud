package presentation

import (
	"time"

	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/policyverificationservice/application/command"
	"github.com/janicaleksander/cloud/policyverificationservice/application/query"
)

//RequestToDomain

func CreatePolicyHTTPToCommand(newID uuid.UUID, p *CreatePolicyRequestDTO) *command.CreatePolicyCommand {
	return &command.CreatePolicyCommand{
		ID:     newID.String(),
		UserID: p.UserID,
		VIN:    p.VIN,
		From:   p.From,
		To:     p.To,
	}
}

func GetPolicyHTTPToQuery(policyID uuid.UUID) *query.GetPolicyQuery {
	return &query.GetPolicyQuery{
		PolicyID: policyID.String(),
	}
}

func GetPoliciesHTTPToQuery() *query.GetPoliciesQuery {
	return &query.GetPoliciesQuery{}
}

func UpdatePolicyTTPToCommand(policyID uuid.UUID, from, to time.Time) *command.UpdatePolicyCommand {
	return &command.UpdatePolicyCommand{
		PolicyID: policyID.String(),
		NewFrom:  from,
		NewTo:    to,
	}
}

func DeletePolicyHTTPToCommand(policyID uuid.UUID) *command.DeletePolicyCommand {
	return &command.DeletePolicyCommand{
		PolicyID: policyID.String(),
	}
}
