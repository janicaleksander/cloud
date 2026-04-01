package persistance

import (
	"context"

	"github.com/janicaleksander/cloud/decisionservice/domain"
	"gorm.io/gorm"
)

type DecisionRepository struct {
	gorm *gorm.DB
}

func NewDecisionRepository(gorm *gorm.DB) *DecisionRepository {
	return &DecisionRepository{gorm: gorm}
}

func (d *DecisionRepository) Save(decision *domain.Decision) (*domain.Decision, error) {
	decisionModel := DomainToDecisionModel(decision)
	err := gorm.G[DecisionModel](d.gorm).Create(context.Background(), decisionModel)
	return DecisionModelToDomain(decisionModel), err

}

// todo check repositories function, especially GetBYid...
// todo check the get generic queries
func (d *DecisionRepository) GetByID(decision uint) (*domain.Decision, error) {
	decisionModel, err := gorm.G[DecisionModel](d.gorm).Where("id = ? ", decision).First(context.Background())
	if err != nil {
		return nil, err
	}
	return DecisionModelToDomain(&decisionModel), nil

}

func (d *DecisionRepository) GetAll() ([]*domain.Decision, error) {
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
	decisionModel := DomainToDecisionModel(decision)

	if err := d.gorm.Save(decisionModel).Error; err != nil {
		return nil, err
	}
	return DecisionModelToDomain(decisionModel), nil
}

func (d *DecisionRepository) DeleteById(id uint) error {
	_, err := gorm.G[DecisionModel](d.gorm).Where("id = ?", id).Delete(context.Background())
	return err
}
