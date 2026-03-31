package persistance

import (
	"context"

	"github.com/janicaleksander/cloud/claimservice/domain"
	"gorm.io/gorm"
)

type ClaimRepository struct {
	gorm *gorm.DB
}

func NewClaimRepository(g *gorm.DB) *ClaimRepository {
	return &ClaimRepository{gorm: g}
}

func (r *ClaimRepository) Save(ctx context.Context, c *domain.Claim) (*domain.Claim, error) {
	claimModel, err := ClaimDomainToModel(c)
	if err != nil {
		return nil, err
	}
	err = gorm.G[ClaimModel](r.gorm).Create(ctx, claimModel)
	if err != nil {
		return nil, err
	}
	claimDomain, err := ClaimModelToDomain(claimModel)
	if err != nil {
		return nil, err
	}
	return claimDomain, nil
}

func (r *ClaimRepository) GetAll(ctx context.Context) ([]*domain.Claim, error) {
	claimModels, err := gorm.G[ClaimModel](r.gorm).Preload("Files", nil).Find(ctx)
	if err != nil {
		return nil, err
	}
	claimDomains := make([]*domain.Claim, 0, 32)
	for idx := range claimModels {
		domainClaim, err := ClaimModelToDomain(&claimModels[idx])
		if err != nil {
			return nil, err
		}
		claimDomains = append(claimDomains, domainClaim)
	}
	return claimDomains, nil
}

func (r *ClaimRepository) GetById(ctx context.Context, id uint) (*domain.Claim, error) {
	claimModel, err := gorm.G[ClaimModel](r.gorm).Preload("Files", nil).Where("id = ?", id).First(ctx)
	if err != nil {
		return nil, err
	}
	claimDomain, err := ClaimModelToDomain(&claimModel)
	if err != nil {
		return nil, err
	}
	return claimDomain, nil
}

func (r *ClaimRepository) Update(ctx context.Context, c *domain.Claim) (*domain.Claim, error) {
	claimModel, err := ClaimDomainToModel(c)
	if err != nil {
		return nil, err
	}
	_, err = gorm.G[ClaimModel](r.gorm).Preload("Files", nil).Where("id = ?", claimModel.ID).Updates(ctx, *claimModel)
	if err != nil {
		return nil, err
	}
	claimDomain, err := ClaimModelToDomain(claimModel)
	if err != nil {
		return nil, err
	}
	return claimDomain, nil
}

func (r *ClaimRepository) DeleteById(ctx context.Context, id uint) error {
	_, err := gorm.G[ClaimModel](r.gorm).Preload("Files", nil).Where("id = ?", id).Delete(ctx)
	return err
}
