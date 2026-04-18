package presentation

import (
	"os"

	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/claimservice/application/command"
	"github.com/janicaleksander/cloud/claimservice/application/query"
	"github.com/janicaleksander/cloud/claimservice/domain"
)

func HTTPCreateClaimRequestToDomain(dto *CreateClaimRequestDTO) *domain.Claim {
	id, _ := uuid.Parse(dto.UserID)
	return &domain.Claim{
		UserID:       id,
		AccidentDate: dto.AccidentDate,
		Email:        dto.Email,
		VIN:          dto.VIN,
	}
}

func HTTPGetClaimDomainToResponse(claim *domain.Claim) *GetClaimResponseDTO {
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

func CreateClaimRequestHTTPToCommand(dto *CreateClaimRequestDTO, files []*os.File) *command.CreateClaimCommand {
	return &command.CreateClaimCommand{
		ID:           uuid.New().String(),
		UserID:       dto.UserID,
		Email:        dto.Email,
		VIN:          dto.VIN,
		AccidentDate: dto.AccidentDate,
		ObjectFiles:  files,
	}
}

func GetClaimRequestHTTPToQuery(claimID string) *query.GetClaimByIdQuery {
	return &query.GetClaimByIdQuery{ClaimID: claimID}
}

func GetClaimsRequestHTTPToQuery() *query.GetClaimsQuery {
	return &query.GetClaimsQuery{}
}

func DeleteClaimRequestHTTPToCommand(claimID string) *command.DeleteClaimCommand {
	return &command.DeleteClaimCommand{ClaimID: claimID}
}

func UpdateClaimRequestHTTPToCommand(claimID string, dto *UpdateClaimRequestDTO) *command.UpdateClaimCommand {
	return &command.UpdateClaimCommand{
		ClaimID:  claimID,
		NewEmail: dto.Email,
	}
}

func GetFileRequestHTTPToQuery(fileID string) *query.GetFileFromStorageQuery {
	return &query.GetFileFromStorageQuery{FileID: fileID}
}
