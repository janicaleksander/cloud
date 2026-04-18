package persistence

import (
	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/claimservice/domain"
)

func ClaimModelToDomain(c *ClaimModel) (*domain.Claim, error) {
	domainFiles := make([]*domain.File, 0, len(c.Files))
	for idx := range c.Files {
		domainFiles = append(domainFiles, &domain.File{
			ID:         c.Files[idx].ID,
			FileName:   c.Files[idx].FileName,
			FileExt:    c.Files[idx].FileExt,
			FileSize:   c.Files[idx].FileSize,
			UploadedAt: c.Files[idx].UploadedAt,
			StorageURL: c.Files[idx].StorageURL,
		})
	}
	status, err := domain.StringToStatus(c.Status)
	if err != nil {
		return nil, err
	}
	claimDomain := &domain.Claim{
		ID:           c.ID,
		UserID:       c.UserID,
		Email:        c.Email,
		VIN:          c.VIN,
		AccidentDate: c.AccidentDate,
		Status:       status,
		Files:        domainFiles,
		UpdatedAt:    c.UpdatedAt,
	}
	return claimDomain, nil

}

func ClaimDomainToModel(c *domain.Claim) (*ClaimModel, error) {
	modelFiles := make([]FileModel, 0, len(c.Files))
	for idx := range c.Files {
		modelFiles = append(modelFiles, FileModel{
			ID:           c.Files[idx].ID,
			FileName:     c.Files[idx].FileName,
			FileExt:      c.Files[idx].FileExt,
			FileSize:     c.Files[idx].FileSize,
			UploadedAt:   c.Files[idx].UploadedAt,
			StorageURL:   c.Files[idx].StorageURL,
			ClaimModelID: c.ID,
		})
	}
	var claimModel = &ClaimModel{
		ID:           c.ID,
		UserID:       c.UserID,
		VIN:          c.VIN,
		Email:        c.Email,
		AccidentDate: c.AccidentDate,
		Status:       string(c.Status),
		Files:        modelFiles,
		UpdatedAt:    c.UpdatedAt,
	}
	return claimModel, nil

}

func FileModelToDomain(f *FileModel) *domain.File {
	return &domain.File{
		ID:         f.ID,
		FileName:   f.FileName,
		FileExt:    f.FileExt,
		FileSize:   f.FileSize,
		UploadedAt: f.UploadedAt,
		StorageURL: f.StorageURL,
	}
}

func FileDomainToModel(f *domain.File, claimID uuid.UUID) *FileModel {
	return &FileModel{
		ID:           f.ID,
		FileName:     f.FileName,
		FileExt:      f.FileExt,
		FileSize:     f.FileSize,
		UploadedAt:   f.UploadedAt,
		StorageURL:   f.StorageURL,
		ClaimModelID: claimID,
	}
}
