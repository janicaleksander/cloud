package presentation

import "github.com/janicaleksander/cloud/claimservice/domain"

func CreateClaimRequestToDomain(dto *CreateClaimRequestDTO) *domain.Claim {
	return &domain.Claim{
		UserID:       dto.UserID,
		CarID:        dto.CarID,
		AccidentDate: dto.AccidentDate,
		VIN:          dto.VIN,
	}
}

func GetClaimDomainToResponse(claim *domain.Claim) *GetClaimResponseDTO {
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
		ID:           claim.ID,
		UserID:       claim.UserID,
		CarID:        claim.CarID,
		AccidentDate: claim.AccidentDate,
		VIN:          claim.VIN,
		Status:       string(claim.Status),
		Files:        files,
		CreatedAt:    claim.CreatedAt,
		UpdatedAt:    claim.UpdatedAt,
	}

}
