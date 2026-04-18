package query

import (
	"io"

	"github.com/janicaleksander/cloud/claimservice/domain"
)

func GetClaimQueryResponseToDomain(query *GetClaimByIdQueryResponse) *domain.Claim {
	newStatus, err := domain.StringToStatus(query.Status)
	if err != nil {
		newStatus = ""
	}
	files := make([]*domain.File, len(query.Files))
	for i, file := range query.Files {
		files[i] = &domain.File{
			ID:         file.ID,
			FileName:   file.FileName,
			FileExt:    file.FileExt,
			FileSize:   file.FileSize,
			UploadedAt: file.UploadedAt,
			StorageURL: file.StorageURL,
		}
	}
	return &domain.Claim{
		ID:           query.ID,
		UserID:       query.UserID,
		Email:        query.Email,
		VIN:          query.VIN,
		AccidentDate: query.AccidentDate,
		Status:       newStatus,
		Files:        files,
		UpdatedAt:    query.UpdatedAt,
	}
}

func GetClaimsQueryResponseToDomain(query []*GetClaimByIdQueryResponse) []*domain.Claim {
	claims := make([]*domain.Claim, len(query))
	for i, claim := range query {
		claims[i] = GetClaimQueryResponseToDomain(claim)
	}
	return claims
}
func GetFileFromStorageQueryResponseToDomain(query *GetFileFromStorageQueryResponse) io.ReadCloser {
	return query.reader
}
