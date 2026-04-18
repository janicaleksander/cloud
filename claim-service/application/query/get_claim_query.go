package query

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/claimservice/domain"
	"github.com/mehdihadeli/go-mediatr"
)

type GetClaimByIdQuery struct {
	ClaimID string `json:"claim_id"`
}

type GetClaimByIdQueryResponse struct {
	ID           string
	UserID       string
	Email        string
	VIN          string
	AccidentDate time.Time
	Status       string
	Files        []FileResponse
	UpdatedAt    time.Time
}

type FileResponse struct {
	ID         string
	FileName   string
	FileExt    string
	FileSize   int64
	UploadedAt time.Time
	StorageURL string
}

type GetClaimQueryHandler struct {
	repo domain.ClaimRepository
}

func NewGetClaimQueryHandler(repo domain.ClaimRepository) *GetClaimQueryHandler {
	return &GetClaimQueryHandler{repo: repo}
}
func (h *GetClaimQueryHandler) SelfRegister() error {
	return mediatr.RegisterRequestHandler[*GetClaimByIdQuery, *GetClaimByIdQueryResponse](h)
}

func (h *GetClaimQueryHandler) Handle(ctx context.Context, query *GetClaimByIdQuery) (*GetClaimByIdQueryResponse, error) {
	qid, err := uuid.Parse(query.ClaimID)
	if err != nil {
		return nil, err
	}
	claimDomain, err := h.repo.GetById(ctx, qid)
	if err != nil {
		return nil, err
	}
	filesResponse := make([]FileResponse, len(claimDomain.Files))
	for i, file := range claimDomain.Files {
		filesResponse[i] = FileResponse{
			ID:         file.ID.String(),
			FileName:   file.FileName,
			FileExt:    file.FileExt,
			FileSize:   file.FileSize,
			UploadedAt: file.UploadedAt,
			StorageURL: file.StorageURL,
		}
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
	}, err

}
