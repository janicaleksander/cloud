package event

type ClaimSubmittedEvent struct {
	UserID  int `json:"userID"`
	ClaimID int `json:"claimID"`
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
