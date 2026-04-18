package presentation

import (
	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/claimservice/domain"
)

func CreateClaimRequestToDomain(dto *CreateClaimRequestDTO) *domain.Claim {
	id, _ := uuid.FromBytes([]byte(dto.UserID))
	return &domain.Claim{
		UserID:       id,
		AccidentDate: dto.AccidentDate,
		Email:        dto.Email,
		VIN:          dto.VIN,
	}
}

func GetClaimDomainToResponse(claim *domain.Claim) *GetClaimResponseDTO {
	files := make([]FileResponseDTO, 0, len(claim.Files))

	for _, f := range claim.Files {
		files = append(files, FileResponseDTO{
			ID:         f.ID.String(),
			FileName:   f.FileName,
			FileExt:    f.FileExt,
			FileSize:   f.FileSize,
			UploadedAt: f.UploadedAt,
			StorageURL: f.StorageURL,
		})
	}
	return &GetClaimResponseDTO{
		ID:           claim.ID.String(),
		UserID:       claim.UserID.String(),
		AccidentDate: claim.AccidentDate,
		VIN:          claim.VIN,
		Email:        claim.Email,
		Status:       string(claim.Status),
		Files:        files,
		UpdatedAt:    claim.UpdatedAt,
	}

}
