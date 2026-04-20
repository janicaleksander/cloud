package persistence

import (
	"github.com/janicaleksander/cloud/valuationservice/domain"
)

func ValuationModelToDomain(m *ValuationModel) *domain.Valuation {
	parts := make([]*domain.Part, 0)
	for _, part := range m.Parts {
		parts = append(parts, &domain.Part{
			ID:   part.ID,
			Name: part.Name,
			Cost: part.Cost,
		})
	}
	return &domain.Valuation{
		ID:      m.ID,
		ClaimID: m.ClaimID,
		Amount:  m.Amount,
		Parts:   parts,
	}
}

func ValuationDomainToModel(d *domain.Valuation) *ValuationModel {
	parts := make([]PartModel, 0)
	for _, part := range d.Parts {
		parts = append(parts, PartModel{
			ID:   part.ID,
			Name: part.Name,
			Cost: part.Cost,
		})

	}
	return &ValuationModel{
		ID:      d.ID,
		ClaimID: d.ClaimID,
		Amount:  d.Amount,
		Parts:   parts,
	}
}
