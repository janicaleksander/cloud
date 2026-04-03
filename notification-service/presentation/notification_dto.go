package presentation

import "time"

type GetNotificationResponseDTO struct {
	ID      uint      `json:"id"`
	ClaimID uint      `json:"claim_id"`
	Body    string    `json:"body"`
	SentTo  string    `json:"sent_to"`
	Time    time.Time `json:"time"`
}
