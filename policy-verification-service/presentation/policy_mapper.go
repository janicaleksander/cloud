package presentation

import (
	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/policyverificationservice/application/command"
	"github.com/janicaleksander/cloud/policyverificationservice/domain"
)

//RequestToDomain

func CreatePolicyRequestToDomain(r *CreatePolicyRequestDTO) *domain.Policy {
	uid, err := uuid.Parse(r.UserID)
	if err != nil {
		return nil
	}
	return &domain.Policy{
		UserID: uid,
		VIN:    r.VIN,
		From:   r.From,
		To:     r.To,
	}
}

// domain to response

func GetPolicyDomainToResponse(p *domain.Policy) *GetPolicyResponseDTO {
	return &GetPolicyResponseDTO{
		ID:     p.ID.String(),
		UserID: p.UserID.String(),
		VIN:    p.VIN,
		From:   p.From,
		To:     p.To,
	}
}

func CreatePolicyResponseHTTPToCommand(newID uuid.UUID, p *CreatePolicyRequestDTO) *command.CreatePolicyCommand {
	uid, err := uuid.Parse(p.UserID)
	if err != nil {
		return nil
	}
	return &command.CreatePolicyCommand{
		ID:     newID,
		UserID: uid,
		VIN:    p.VIN,
		From:   p.From,
		To:     p.To,
	}
}
