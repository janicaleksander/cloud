package api

import (
	"github.com/janicaleksander/cloud/policyverificationservice/domain"
)

//RequestToDomain

func CreatePolicyRequestToDomain(r *CreatePolicyRequestDTO) *domain.Policy {
	return &domain.Policy{
		UserID: r.UserID,
		VIN:    r.VIN,
		From:   r.From,
		To:     r.To,
	}
}

// domain to response

func GetPolicyDomainToResponse(p *domain.Policy) *GetPolicyResponseDTO {
	return &GetPolicyResponseDTO{
		ID:     p.ID,
		UserID: p.UserID,
		VIN:    p.VIN,
		From:   p.From,
		To:     p.To,
	}
}
