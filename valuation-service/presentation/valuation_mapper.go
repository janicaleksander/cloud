package presentation

import (
	"github.com/janicaleksander/cloud/valuationservice/domain"
)

func GetValuationDomainToResponse(v *domain.Valuation) *GetValuationResponseDTO {
	parts := make([]Part, len(v.Parts))
	for i, part := range v.Parts {
		parts[i] = Part{
			ID:   part.ID,
			Name: part.Name,
			Cost: part.Cost,
		}
	}
	return &GetValuationResponseDTO{
		ID:      v.ID,
		ClaimID: v.ClaimID,
		Amount:  v.Amount,
		Parts:   parts,
	}
}
