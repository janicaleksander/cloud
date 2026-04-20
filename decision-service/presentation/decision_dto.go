package presentation

type GetDecisionResponseDTO struct {
	ID         string  `json:"id"`
	EmployeeID string  `json:"employee_id,omitempty"`
	ClaimID    string  `json:"claim_id"`
	Payout     float64 `json:"payout"`
	State      string  `json:"state"`
}

type UpdateDecisionRequestDTO struct {
	EmpID    string `json:"emp_id"`
	NewState string `json:"new_status"`
	Reason   string `json:"reason,omitempty"`
}
