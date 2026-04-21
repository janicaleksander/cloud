package query

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/claimservice/domain"
	"github.com/mehdihadeli/go-mediatr"
)

type GetClaimByIdQuery struct {
	ClaimID string
}

type GetClaimByIdQueryResponse struct {
	ID           string         `json:"id"`
	UserID       string         `json:"user_id"`
	Email        string         `json:"email"`
	VIN          string         `json:"vin"`
	AccidentDate time.Time      `json:"accident_date"`
	Status       string         `json:"status"`
	Files        []FileResponse `json:"files"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type FileResponse struct {
	ID         string    `json:"id"`
	FileName   string    `json:"file_name"`
	FileExt    string    `json:"file_ext"`
	FileSize   int64     `json:"file_size"`
	UploadedAt time.Time `json:"uploaded_at"`
	StorageURL string    `json:"storage_url"`
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
		filesResponse[i] = *FileDomainToQueryResponse(file)
	}
	return ClaimDomainToQueryResponse(claimDomain), err

}
