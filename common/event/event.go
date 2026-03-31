package event

import "time"

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
	ClaimID      uint         `json:"claim_id"`
	PayoutAmount float64      `json:"payout_amount"`
	DamageItems  []DamageItem `json:"damage_items"`
}
type DamageItem struct {
	Part     string  `json:"part"`
	Severity string  `json:"severity"`
	Cost     float64 `json:"cost"`
}

type PayoutApprovedEvent struct {
	ClaimID      uint    `json:"claim_id"`
	PayoutAmount float64 `json:"payout_amount"`
}

type PayoutRejectedEvent struct {
	ClaimID uint   `json:"claim_id"`
	Reason  string `json:"reason"`
}
