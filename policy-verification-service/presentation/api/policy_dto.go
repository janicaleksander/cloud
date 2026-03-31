package api

import "time"

type CreatePolicyRequestDTO struct {
	UserID uint      `json:"user_id"`
	VIN    string    `json:"vin"`
	From   time.Time `json:"from"`
	To     time.Time `json:"to"`
}

type GetPolicyResponseDTO struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	VIN       string    `json:"vin"`
	From      time.Time `json:"from"`
	To        time.Time `json:"to"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdatePolicyRequest struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to,omitempty"`
}
