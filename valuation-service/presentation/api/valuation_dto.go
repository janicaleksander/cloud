package presentation

import "time"

type GetValuationResponseDTO struct {
	ID        uint      `json:"id"`
	ClaimID   uint      `json:"claim_id"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateValuationRequestDTO struct {
	Amount float64
}
