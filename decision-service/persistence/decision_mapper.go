package persistence

import (
	"github.com/janicaleksander/cloud/decisionservice/domain"
	"gorm.io/gorm"
)

func DecisionModelToDomain(decision *DecisionModel) *domain.Decision {
	return &domain.Decision{
		ID:         decision.ID,
		ClaimID:    decision.ClaimID,
		EmployeeID: decision.EmployeeID,
		Result:     domain.StringToResult(decision.Result),
		Payout:     decision.Payout,
	}
}

func DomainToDecisionModel(decision *domain.Decision) *DecisionModel {
	return &DecisionModel{
		Model:      gorm.Model{ID: decision.ID},
		ClaimID:    decision.ClaimID,
		EmployeeID: decision.EmployeeID,
		Payout:     decision.Payout,
		Result:     string(decision.Result),
	}
}
