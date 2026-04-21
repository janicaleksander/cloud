package query

import "github.com/janicaleksander/cloud/decisionservice/domain"

func DecisionDomainToQueryResponse(d *domain.Decision) *GetDecisionQueryResult {
	return &GetDecisionQueryResult{
		ID:         d.ID.String(),
		EmployeeID: d.EmployeeID.String(),
		ClaimID:    d.ClaimID.String(),
		Payout:     d.Payout,
		State:      string(d.Result),
	}
}
