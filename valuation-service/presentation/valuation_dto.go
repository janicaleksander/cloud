package presentation

type GetValuationResponseDTO struct {
	ID      string  `json:"id"`
	ClaimID string  `json:"claim_id"`
	Amount  float64 `json:"amount"`
	Parts   []Part  `json:"parts"`
}

type Part struct {
	ID   string  `json:"id"`
	Name string  `json:"name"`
	Cost float64 `json:"cost"`
}

/*
type UpdateValuationRequestDTO struct {
	Amount float64
}

*/
