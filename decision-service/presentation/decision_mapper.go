package presentation

import "github.com/janicaleksander/cloud/decisionservice/domain"

func GetDecisionDomainToResponse(d *domain.Decision) *GetDecisionResponseDTO {
	return &GetDecisionResponseDTO{
		ID:         d.ID,
		EmployeeID: d.EmployeeID,
		ClaimID:    d.ClaimID,
		Payout:     d.Payout,
		State:      string(d.Result),
	}
}
