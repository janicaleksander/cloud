package query

import (
	"log/slog"

	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/claimservice/domain"
)

func GetClaimQueryResponseToDomain(query *GetClaimByIdQueryResponse) *domain.Claim {
	newStatus, err := domain.StringToStatus(query.Status)
	if err != nil {
		newStatus = ""
	}
	files := make([]*domain.File, len(query.Files))
	for i, file := range query.Files {
		fid, err := uuid.Parse(file.ID)
		if err != nil {
			slog.Error("Error converting file ID to UUID", "error", err)
			continue
		}
		files[i] = &domain.File{
			ID:         fid,
			FileName:   file.FileName,
			FileExt:    file.FileExt,
			FileSize:   file.FileSize,
			UploadedAt: file.UploadedAt,
			StorageURL: file.StorageURL,
		}
	}
	qid, err := uuid.Parse(query.ID)
	if err != nil {
		slog.Error("Error converting claim ID to UUID", "error", err)
		return nil
	}
	uid, err := uuid.Parse(query.UserID)
	if err != nil {
		slog.Error("Error converting user ID to UUID", "error", err)
		return nil
	}
	return &domain.Claim{
		ID:           qid,
		UserID:       uid,
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
