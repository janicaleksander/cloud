package persistance

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/policyverificationservice/domain"
	"gorm.io/gorm"
)

type PolicyRepository struct {
	gorm *gorm.DB
}

func NewPolicyRepository(gorm *gorm.DB) *PolicyRepository {
	slog.Info("Initializing PolicyRepository")
	return &PolicyRepository{gorm: gorm}
}

func (pr *PolicyRepository) GetAll(ctx context.Context) ([]*domain.Policy, error) {
	slog.Info("Getting all policies from the database")
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

func (pr *PolicyRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.Policy, error) {
	slog.Info("Getting policy by ID from the database", "id", id)
	policyModel, err := gorm.G[PolicyModel](pr.gorm).Where("id = ?", id).First(ctx)
	if err != nil {
		return nil, err
	}
	return PolicyModelToDomain(&policyModel), nil
}

func (pr *PolicyRepository) Save(ctx context.Context, p *domain.Policy) (*domain.Policy, error) {
	slog.Info("Saving policy to the database")
	policyModel := PolicyDomainToModel(p)
	err := gorm.G[PolicyModel](pr.gorm).Create(ctx, policyModel)
	if err != nil {
		return nil, err
	}
	return PolicyModelToDomain(policyModel), nil
}

func (pr *PolicyRepository) Update(ctx context.Context, p *domain.Policy) (*domain.Policy, error) {
	slog.Info("Updating policy in the database", "policy", p)
	policyModel := PolicyDomainToModel(p)

	if err := pr.gorm.WithContext(ctx).Save(policyModel).Error; err != nil {
		return nil, err
	}

	return PolicyModelToDomain(policyModel), nil
}
func (pr *PolicyRepository) DeleteById(ctx context.Context, id uuid.UUID) error {
	slog.Info("Deleting policy by ID from the database", "id", id)
	_, err := gorm.G[PolicyModel](pr.gorm).Where("id = ?", id).Delete(ctx)
	return err

}
func (pr *PolicyRepository) IfUserHasPolicy(ctx context.Context, userID uuid.UUID, vin string) (bool, *domain.Policy) {
	slog.Info("Checking if user has policy for given VIN", "userID", userID, "vin", vin)
	p, err := gorm.G[PolicyModel](pr.gorm).
		Where("user_id = ? AND vin = ?", userID, vin).
		First(ctx)

	if err != nil {
		return false, nil
	}

	return true, PolicyModelToDomain(&p)
}
