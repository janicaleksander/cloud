package persistance

import (
	"context"

	"github.com/janicaleksander/cloud/valuationservice/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ValuationRepository struct {
	gorm *gorm.DB
}

func NewValuationRepository(gorm *gorm.DB) *ValuationRepository {
	return &ValuationRepository{gorm: gorm}
}

func (v *ValuationRepository) GetAll(ctx context.Context) ([]*domain.Valuation, error) {
	valuationModels, err := gorm.G[ValuationModel](v.gorm).Find(ctx)
	if err != nil {
		return nil, err
	}
	valuationDomains := make([]*domain.Valuation, 0, len(valuationModels))
	for idx := range valuationModels {
		domainValuation := ValuationModelToDomain(&valuationModels[idx])
		valuationDomains = append(valuationDomains, domainValuation)
	}
	return valuationDomains, nil
}
func (v *ValuationRepository) GetById(ctx context.Context, id uint) (*domain.Valuation, error) {
	valuationModel, err := gorm.G[ValuationModel](v.gorm).Where("id = ?", id).First(ctx)
	if err != nil {
		return nil, err
	}
	domainValuation := ValuationModelToDomain(&valuationModel)
	return domainValuation, nil
}
func (v *ValuationRepository) Save(ctx context.Context, domainValuation *domain.Valuation) (*domain.Valuation, error) {
	domainModel := ValuationDomainToModel(domainValuation)
	err := gorm.G[ValuationModel](v.gorm).Create(ctx, domainModel)
	if err != nil {
		return nil, err
	}
	return ValuationModelToDomain(domainModel), nil

}
func (v *ValuationRepository) Update(ctx context.Context, valuationDomain *domain.Valuation) (*domain.Valuation, error) {
	valuationModel := ValuationDomainToModel(valuationDomain)
	var updated ValuationModel

	err := v.gorm.
		Model(&ValuationModel{}).
		Where("id = ?", valuationModel.ID).
		Clauses(clause.Returning{}).
		Updates(valuationModel).
		Scan(&updated).Error

	if err != nil {
		return nil, err
	}
	valuationDomainn := ValuationModelToDomain(&updated)
	return valuationDomainn, nil
}
func (v *ValuationRepository) DeleteById(ctx context.Context, id uint) error {
	_, err := gorm.G[ValuationModel](v.gorm).Where("id = ?", id).Delete(ctx)
	return err

}
