package domain

type ClaimSubmittedEvent struct {
	UserID  int `json:"userID"`
	ClaimID int `json:"claimID"`
}
