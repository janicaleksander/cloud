package event

import "time"

type RegisterUserForNotificationEvent struct {
	ClaimID string `json:"claim_id"`
	Email   string `json:"email"`
}
type ChangeEmailForNotification struct {
	ClaimID string `json:"claim_id"`
	Email   string `json:"email"`
}
type ClaimSubmittedEvent struct {
	ClaimID      string    `json:"claim_id"`
	UserID       string    `json:"user_id"`
	VIN          string    `json:"vin"`
	AccidentDate time.Time `json:"accident_date"`
	StorageURL   []string  `json:"storage_url"`
}

type PolicyVerifiedEvent struct {
	ClaimID    string   `json:"claim_id"`
	StorageURL []string `json:"storage_url"`
}

type PolicyDeniedEvent struct {
	ClaimID string `json:"claim_id"`
	Reason  string `json:"reason"`
}

type ValuationCalculatedEvent struct {
	ClaimID      string  `json:"claim_id"`
	PayoutAmount float64 `json:"payout_amount"`
}

type PayoutApprovedEvent struct {
	ClaimID              string  `json:"claim_id"`
	ByEmployeeID         string  `json:"by_employee_id"`
	AcceptedPayoutAmount float64 `json:"payout_amount"`
}

type PayoutRejectedEvent struct {
	ClaimID      string `json:"claim_id"`
	ByEmployeeID string `json:"by_employee_id"`
	Reason       string `json:"reason"`
}
