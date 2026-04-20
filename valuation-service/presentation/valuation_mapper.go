package presentation

import (
	"github.com/janicaleksander/cloud/valuationservice/domain"
)

func GetValuationDomainToResponse(v *domain.Valuation) *GetValuationResponseDTO {
	parts := make([]Part, len(v.Parts))
	for i, part := range v.Parts {
		parts[i] = Part{
			ID:   part.ID.String(),
			Name: part.Name,
			Cost: part.Cost,
		}
	}
	return &GetValuationResponseDTO{
		ID:      v.ID.String(),
		ClaimID: v.ClaimID.String(),
		Amount:  v.Amount,
		Parts:   parts,
	}
}
