package persistance

import (
	"context"
	"log"

	"github.com/janicaleksander/cloud/policyverificationservice/domain"
	"gorm.io/gorm"
)

type PolicyRepository struct {
	gorm *gorm.DB
}

func NewPolicyRepository(gorm *gorm.DB) *PolicyRepository {
	return &PolicyRepository{gorm: gorm}
}

func (pr *PolicyRepository) GetAll(ctx context.Context) ([]*domain.Policy, error) {
	policiesModel, err := gorm.G[PolicyModel](pr.gorm).Find(ctx)
	if err != nil {
		return nil, err
	}
	domainPolicies := make([]*domain.Policy, 0, len(policiesModel))
	for idx := range policiesModel {
		domainPolicies = append(domainPolicies, PolicyModelToDomain(&policiesModel[idx]))
	}
	return domainPolicies, nil

}

func (pr *PolicyRepository) GetById(ctx context.Context, id uint) (*domain.Policy, error) {
	policyModel, err := gorm.G[PolicyModel](pr.gorm).Where("id = ?", id).First(ctx)
	if err != nil {
		return nil, err
	}
	return PolicyModelToDomain(&policyModel), nil
}

func (pr *PolicyRepository) Save(ctx context.Context, p *domain.Policy) (*domain.Policy, error) {
	policyModel := PolicyDomainToModel(p)
	err := gorm.G[PolicyModel](pr.gorm).Create(ctx, policyModel)
	if err != nil {
		return nil, err
	}
	return PolicyModelToDomain(policyModel), nil
}

func (pr *PolicyRepository) Update(ctx context.Context, p *domain.Policy) (*domain.Policy, error) {
	policyModel := PolicyDomainToModel(p)

	if err := pr.gorm.WithContext(ctx).Save(policyModel).Error; err != nil {
		return nil, err
	}

	return PolicyModelToDomain(policyModel), nil
}
func (pr *PolicyRepository) DeleteById(ctx context.Context, id uint) error {
	_, err := gorm.G[PolicyModel](pr.gorm).Where("id = ?", id).Delete(ctx)
	return err

}
func (pr *PolicyRepository) IfUserHasPolicy(ctx context.Context, userID uint, vin string) (bool, *domain.Policy) {
	p, err := gorm.G[PolicyModel](pr.gorm).
		Where("user_id = ? AND vin = ?", userID, vin).
		First(ctx)

	log.Printf("IfUserHasPolicy: userID=%d vin=%s err=%v p=%+v", userID, vin, err, p)
	if err != nil {
		return false, nil
	}

	return true, PolicyModelToDomain(&p)
}
