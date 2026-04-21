package query

import (
	"github.com/janicaleksander/cloud/claimservice/domain"
)

func FileDomainToQueryResponse(fileDomain *domain.File) *FileResponse {
	return &FileResponse{
		ID:         fileDomain.ID.String(),
		FileName:   fileDomain.FileExt,
		FileExt:    fileDomain.FileExt,
		FileSize:   fileDomain.FileSize,
		UploadedAt: fileDomain.UploadedAt,
		StorageURL: fileDomain.StorageURL,
	}
}

func ClaimDomainToQueryResponse(claimDomain *domain.Claim) *GetClaimByIdQueryResponse {
	filesResponse := make([]FileResponse, len(claimDomain.Files))
	for i, file := range claimDomain.Files {
		filesResponse[i] = *FileDomainToQueryResponse(file)
	}
	return &GetClaimByIdQueryResponse{
		ID:           claimDomain.ID.String(),
		UserID:       claimDomain.UserID.String(),
		Email:        claimDomain.Email,
		VIN:          claimDomain.VIN,
		AccidentDate: claimDomain.AccidentDate,
		Status:       string(claimDomain.Status),
		Files:        filesResponse,
		UpdatedAt:    claimDomain.UpdatedAt,
	}
}
