package persistence

import (
	"context"
	"log/slog"

	"github.com/janicaleksander/cloud/valuationservice/domain"
	"gorm.io/gorm"
)

type ValuationRepository struct {
	gorm *gorm.DB
}

func NewValuationRepository(gorm *gorm.DB) *ValuationRepository {
	slog.Info("Initializing ValuationRepository")
	return &ValuationRepository{gorm: gorm}
}

func (v *ValuationRepository) GetAll(ctx context.Context) ([]*domain.Valuation, error) {
	slog.Info("Getting all valuations from the database")
	valuationModels, err := gorm.G[ValuationModel](v.gorm).Preload("Parts", nil).Find(ctx)
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
	slog.Info("Getting valuation by ID from the database", "id", id)
	valuationModel, err := gorm.G[ValuationModel](v.gorm).Preload("Parts", nil).Where("id = ?", id).First(ctx)
	if err != nil {
		return nil, err
	}
	domainValuation := ValuationModelToDomain(&valuationModel)
	return domainValuation, nil
}
func (v *ValuationRepository) Save(ctx context.Context, domainValuation *domain.Valuation) (*domain.Valuation, error) {
	slog.Info("Saving valuation to the database", "valuationID", domainValuation.ID)
	domainModel := ValuationDomainToModel(domainValuation)
	err := gorm.G[ValuationModel](v.gorm).Create(ctx, domainModel)
	if err != nil {
		return nil, err
	}
	return ValuationModelToDomain(domainModel), nil

}
func (v *ValuationRepository) Update(ctx context.Context, valuationDomain *domain.Valuation) (*domain.Valuation, error) {
	slog.Info("Updating valuation in the database", "valuationID", valuationDomain.ID)
	valuationModel := ValuationDomainToModel(valuationDomain)

	err := v.gorm.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(valuationModel).Error; err != nil {
			return err
		}

		if err := tx.Where("valuation_id = ?", valuationModel.ID).
			Delete(&PartModel{}).Error; err != nil {
			return err
		}

		for i := range valuationModel.Parts {
			valuationModel.Parts[i].ID = 0
			valuationModel.Parts[i].ValuationID = valuationModel.ID
		}

		return tx.Create(&valuationModel.Parts).Error
	})

	if err != nil {
		return nil, err
	}

	var updated ValuationModel
	if err := v.gorm.WithContext(ctx).
		Preload("Parts").
		First(&updated, valuationModel.ID).Error; err != nil {
		return nil, err
	}

	return ValuationModelToDomain(&updated), nil
}
func (v *ValuationRepository) DeleteById(ctx context.Context, id uint) error {
	slog.Info("Deleting valuation by ID from the database", "id", id)
	_, err := gorm.G[ValuationModel](v.gorm).Where("id = ?", id).Delete(ctx)
	return err

}
