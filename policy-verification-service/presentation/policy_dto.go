package presentation

import "time"

type CreatePolicyRequestDTO struct {
	UserID string    `json:"user_id"`
	VIN    string    `json:"vin"`
	From   time.Time `json:"from"`
	To     time.Time `json:"to"`
}

type UpdatePolicyRequest struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to,omitempty"`
}
