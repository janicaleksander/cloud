package presentation

type UpdateDecisionRequestDTO struct {
	EmpID    string `json:"emp_id"`
	NewState string `json:"new_status"`
	Reason   string `json:"reason,omitempty"`
}
