package event

type ClaimSubmittedEvent struct {
	UserID     uint     `json:"user_id"`
	ClaimID    uint     `json:"claim_id"`
	StorageURL []string `json:"storage_url"`
} // this claim_service sends

type PolicyVerifiedEvent struct {
}

type PolicyDeniedEvent struct {
}

type ValuationCalculatedEvent struct {
}

type PayoutApprovedEvent struct {
}

type PayoutRejectedEvent struct {
}
