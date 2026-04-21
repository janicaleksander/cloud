package presentation

import "time"

type CreateClaimRequestDTO struct {
	UserID       string    `json:"user_id"`
	Email        string    `json:"email"`
	AccidentDate time.Time `json:"accident_date"`
	VIN          string    `json:"vin"`
}
