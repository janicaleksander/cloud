package event

type ClaimSubmittedEvent struct {
	UserID     uint     `json:"user_id"`
	ClaimID    uint     `json:"claim_id"`
	StorageURL []string `json:"storage_url"`
} // this claim_service sends

type PolicyVerifiedEvent struct {
	ClaimID uint `json:"claim_id"`
}

type PolicyDeniedEvent struct {
	ClaimID uint `json:"claim_id"`
}

type ValuationCalculatedEvent struct {
	ClaimID uint `json:"claim_id"`
}

type PayoutApprovedEvent struct {
	ClaimID uint `json:"claim_id"`
}

type PayoutRejectedEvent struct {
	ClaimID uint `json:"claim_id"`
}
