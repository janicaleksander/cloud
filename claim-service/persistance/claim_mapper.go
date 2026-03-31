package persistance

import (
	"github.com/janicaleksander/cloud/claimservice/domain"
	"gorm.io/gorm"
)

func ClaimModelToDomain(c *ClaimModel) (*domain.Claim, error) {
	domainFiles := make([]*domain.File, 0, len(c.Files))
	for idx := range c.Files {
		domainFiles = append(domainFiles, &domain.File{
			ID:         c.Files[idx].ID,
			FileName:   c.Files[idx].FileName,
			FileExt:    c.Files[idx].FileExt,
			StorageURL: c.Files[idx].StorageURL,
			CreatedAt:  c.Files[idx].CreatedAt,
			UpdatedAt:  c.Files[idx].UpdatedAt,
		},
		)
	}
	status, err := domain.StringToStatus(c.Status)
	if err != nil {
		return nil, err
	}
	claimDomain := &domain.Claim{
		ID:           c.ID,
		UserID:       c.UserID,
		CarID:        c.CarID,
		AccidentDate: c.AccidentDate,
		Status:       status,
		Files:        domainFiles,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
	}
	return claimDomain, nil

}

func ClaimDomainToModel(c *domain.Claim) (*ClaimModel, error) {
	modelFiles := make([]FileModel, 0, len(c.Files))
	for idx := range c.Files {
		modelFiles = append(modelFiles, FileModel{
			Model: gorm.Model{
				ID:        c.Files[idx].ID,
				CreatedAt: c.Files[idx].CreatedAt,
				UpdatedAt: c.Files[idx].UpdatedAt},
			FileName:     c.Files[idx].FileName,
			FileExt:      c.Files[idx].FileExt,
			StorageURL:   c.Files[idx].StorageURL,
			ClaimModelID: c.ID,
		})
	}
	var claimModel = &ClaimModel{
		Model: gorm.Model{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		},
		UserID:       c.UserID,
		CarID:        c.CarID,
		AccidentDate: c.AccidentDate,
		Status:       string(c.Status),
		Files:        modelFiles,
	}
	return claimModel, nil

}
