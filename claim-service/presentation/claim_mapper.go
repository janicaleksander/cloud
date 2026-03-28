package presentation

import "github.com/janicaleksander/cloud/claimservice/domain"

func CreateClaimRequestToDomain(dto *CreateClaimRequestDTO) *domain.Claim {
	domainFiles := make([]*domain.File, 0, len(dto.Files))
	for _, f := range dto.Files {
		domainFiles = append(domainFiles, &domain.File{
			FileName:   f.FileName,
			FileExt:    f.FileExt,
			StorageURL: f.StorageURL,
		})
	}

	return &domain.Claim{
		UserID: dto.UserID,
		CarID:  dto.CarID,
		Files:  domainFiles,
	}
}

func GetClaimDomainToRequest(claim *domain.Claim) *GetClaimResponseDTO {
	files := make([]FileResponseDTO, 0, len(claim.Files))
	for _, f := range claim.Files {
		files = append(files, FileResponseDTO{
			ID:         f.ID,
			FileName:   f.FileName,
			FileExt:    f.FileExt,
			StorageURL: f.StorageURL,
		})
	}
	return &GetClaimResponseDTO{
		ID:        claim.ID,
		UserID:    claim.UserID,
		CarID:     claim.CarID,
		Status:    string(claim.Status),
		Files:     files,
		CreatedAt: claim.CreatedAt,
		UpdatedAt: claim.UpdatedAt,
	}

}
