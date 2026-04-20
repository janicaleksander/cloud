package presentation

import "github.com/janicaleksander/cloud/decisionservice/domain"

func GetDecisionDomainToResponse(d *domain.Decision) *GetDecisionResponseDTO {
	return &GetDecisionResponseDTO{
		ID:         d.ID.String(),
		EmployeeID: d.EmployeeID.String(),
		ClaimID:    d.ClaimID.String(),
		Payout:     d.Payout,
		State:      string(d.Result),
	}
}
