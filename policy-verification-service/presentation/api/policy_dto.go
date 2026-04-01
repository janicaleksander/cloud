package api

import "time"

type CreatePolicyRequestDTO struct {
	UserID uint      `json:"user_id"`
	VIN    string    `json:"vin"`
	From   time.Time `json:"from"`
	To     time.Time `json:"to"`
}

type GetPolicyResponseDTO struct {
	ID     uint      `json:"id"`
	UserID uint      `json:"user_id"`
	VIN    string    `json:"vin"`
	From   time.Time `json:"from"`
	To     time.Time `json:"to"`
}

type UpdatePolicyRequest struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to,omitempty"`
}
