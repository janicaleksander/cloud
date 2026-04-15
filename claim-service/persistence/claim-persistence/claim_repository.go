package claim_persistence

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/janicaleksander/cloud/claimservice/domain"
	"gorm.io/gorm"
)

type ClaimRepository struct {
	gorm *gorm.DB
}

func NewClaimRepository(g *gorm.DB) *ClaimRepository {
	slog.Info("Initializing ClaimRepository")
	return &ClaimRepository{gorm: g}
}

func (r *ClaimRepository) Save(ctx context.Context, c *domain.Claim) (*domain.Claim, error) {
	slog.Info("Saving claim to database")
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
	slog.Info("Getting all claims from database")
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
	slog.Info("Getting claim by ID from database", "claimID", id)
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
	slog.Info("Updating claim in database", "claimID", c.ID)
	claimModel, err := ClaimDomainToModel(c)
	if err != nil {
		return nil, err
	}

	err = r.gorm.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(claimModel).Error; err != nil {
			return err
		}

		if err := tx.Where("claim_model_id = ?", claimModel.ID).
			Delete(&FileModel{}).Error; err != nil {
			return err
		}

		for i := range claimModel.Files {
			claimModel.Files[i].ID = 0
			claimModel.Files[i].ClaimModelID = claimModel.ID
		}

		if len(claimModel.Files) == 0 {
			return nil
		}
		return tx.Create(&claimModel.Files).Error
	})

	if err != nil {
		return nil, err
	}

	var updated ClaimModel
	if err := r.gorm.WithContext(ctx).
		Preload("Files").
		First(&updated, claimModel.ID).Error; err != nil {
		return nil, err
	}

	return ClaimModelToDomain(&updated)
}
func (r *ClaimRepository) DeleteById(ctx context.Context, id uint) error {
	slog.Info("Deleting claim by ID from database", "claimID", id)
	_, err := gorm.G[ClaimModel](r.gorm).Preload("Files", nil).Where("id = ?", id).Delete(ctx)
	return err
}

func (r *ClaimRepository) UpdateStatus(ctx context.Context, claimID uint, newStatus domain.Status) error {
	slog.Info("Updating claim status in database", "claimID", claimID, "newStatus", newStatus)
	result := r.gorm.WithContext(ctx).
		Model(&ClaimModel{}).
		Where("id = ?", claimID).
		Update("status", string(newStatus))

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("claim %d not found", claimID)
	}
	return nil
}

func (r *ClaimRepository) GetFileById(ctx context.Context, fileID uint) (*domain.File, error) {
	fileModel, err := gorm.G[FileModel](r.gorm).Where("id = ?", fileID).First(ctx)
	if err != nil {
		return nil, err
	}
	return FileModelToDomain(&fileModel), nil

}
