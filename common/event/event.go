package event

import "time"

type RegisterUserForNotificationEvent struct {
	ClaimID uint   `json:"claim_id"`
	Email   string `json:"email"`
}
type ChangeEmailForNotification struct {
	ClaimID uint   `json:"claim_id"`
	Email   string `json:"email"`
}
type ClaimSubmittedEvent struct {
	ClaimID      uint      `json:"claim_id"`
	UserID       uint      `json:"user_id"`
	VIN          string    `json:"vin"`
	AccidentDate time.Time `json:"accident_date"`
	StorageURL   []string  `json:"storage_url"`
}

type PolicyVerifiedEvent struct {
	ClaimID    uint     `json:"claim_id"`
	StorageURL []string `json:"storage_url"`
}

type PolicyDeniedEvent struct {
	ClaimID uint   `json:"claim_id"`
	Reason  string `json:"reason"`
}

type ValuationCalculatedEvent struct {
	ClaimID      uint    `json:"claim_id"`
	PayoutAmount float64 `json:"payout_amount"`
}

type PayoutApprovedEvent struct {
	ClaimID              uint    `json:"claim_id"`
	ByEmployeeID         uint    `json:"by_employee_id"`
	AcceptedPayoutAmount float64 `json:"payout_amount"`
}

type PayoutRejectedEvent struct {
	ClaimID      uint   `json:"claim_id"`
	ByEmployeeID uint   `json:"by_employee_id"`
	Reason       string `json:"reason"`
}
