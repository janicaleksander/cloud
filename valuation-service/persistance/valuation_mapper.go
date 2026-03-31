package persistance

import (
	"github.com/janicaleksander/cloud/valuationservice/domain"
	"gorm.io/gorm"
)

func ValuationModelToDomain(m *ValuationModel) *domain.Valuation {
	return &domain.Valuation{
		ID:        m.ID,
		ClaimID:   m.ClaimID,
		Amount:    m.Amount,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func ValuationDomainToModel(d *domain.Valuation) *ValuationModel {
	return &ValuationModel{
		Model: gorm.Model{
			ID:        d.ID,
			CreatedAt: d.CreatedAt,
			UpdatedAt: d.UpdatedAt,
		},
		ClaimID: d.ClaimID,
		Amount:  d.Amount,
	}
}
