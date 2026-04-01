package presentation

import (
	domain "github.com/janicaleksander/cloud/decisionservice/domain"
)

type GetDecisionResponseDTO struct {
	ID         uint    `json:"id"`
	EmployeeID *uint   `json:"employee_id,omitempty"`
	ClaimID    uint    `json:"claim_id"`
	Payout     float64 `json:"payout"`
	State      string  `json:"state"`
}

type UpdateDecisionRequestDTO struct {
	EmpID    uint                  `json:"emp_id"`
	NewState domain.DecisionResult `json:"new_status"`
	Reason   string                `json:"reason,omitempty"`
}
