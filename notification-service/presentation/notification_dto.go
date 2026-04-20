package presentation

import "time"

type GetNotificationResponseDTO struct {
	ID      string    `json:"id"`
	ClaimID string    `json:"claim_id"`
	Body    string    `json:"body"`
	SentTo  string    `json:"sent_to"`
	Time    time.Time `json:"time"`
}
