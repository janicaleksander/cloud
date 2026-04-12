package persistence

import (
	"context"
	"log/slog"

	"github.com/janicaleksander/cloud/decisionservice/domain"
	"gorm.io/gorm"
)

type DecisionRepository struct {
	gorm *gorm.DB
}

func NewDecisionRepository(gorm *gorm.DB) *DecisionRepository {
	slog.Info("Initializing DecisionRepository")
	return &DecisionRepository{gorm: gorm}
}

func (d *DecisionRepository) Save(decision *domain.Decision) (*domain.Decision, error) {
	slog.Info("Saving decision to database")
	decisionModel := DomainToDecisionModel(decision)
	err := gorm.G[DecisionModel](d.gorm).Create(context.Background(), decisionModel)
	return DecisionModelToDomain(decisionModel), err

}

func (d *DecisionRepository) GetByID(decision uint) (*domain.Decision, error) {
	slog.Info("Getting decision by ID from database", "decisionID", decision)
	decisionModel, err := gorm.G[DecisionModel](d.gorm).Where("id = ? ", decision).First(context.Background())
	if err != nil {
		return nil, err
	}
	return DecisionModelToDomain(&decisionModel), nil

}

func (d *DecisionRepository) GetAll() ([]*domain.Decision, error) {
	slog.Info("Getting all decisions from database")
	decisionModels, err := gorm.G[DecisionModel](d.gorm).Find(context.Background())
	if err != nil {
		return nil, err
	}
	decisions := make([]*domain.Decision, len(decisionModels))
	for i, model := range decisionModels {
		decisions[i] = DecisionModelToDomain(&model)
	}
	return decisions, nil

}
func (d *DecisionRepository) GetAllWaiting() ([]*domain.Decision, error) {
	slog.Info("Getting all waiting decisions from database")
	decisionModels, err := gorm.G[DecisionModel](d.gorm).Where("result = ?", string(domain.WAITING)).Find(context.Background())
	if err != nil {
		return nil, err
	}
	decisions := make([]*domain.Decision, len(decisionModels))
	for i, model := range decisionModels {
		decisions[i] = DecisionModelToDomain(&model)
	}
	return decisions, nil

}

func (d *DecisionRepository) Update(decision *domain.Decision) (*domain.Decision, error) {
	slog.Info("Updating decision in database", "decisionID", decision.ID)
	decisionModel := DomainToDecisionModel(decision)
	if err := d.gorm.Save(decisionModel).Error; err != nil {
		return nil, err
	}
	return DecisionModelToDomain(decisionModel), nil
}

func (d *DecisionRepository) DeleteById(id uint) error {
	slog.Info("Deleting decision by ID from database", "decisionID", id)
	_, err := gorm.G[DecisionModel](d.gorm).Where("id = ?", id).Delete(context.Background())
	return err
}
