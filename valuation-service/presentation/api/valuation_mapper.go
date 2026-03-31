package presentation

import (
	"github.com/janicaleksander/cloud/valuationservice/domain"
)

func GetValuationDomainToResponse(v *domain.Valuation) *GetValuationResponseDTO {
	return &GetValuationResponseDTO{
		ID:        v.ID,
		ClaimID:   v.ClaimID,
		Amount:    v.Amount,
		CreatedAt: v.CreatedAt,
		UpdatedAt: v.UpdatedAt,
	}
}
