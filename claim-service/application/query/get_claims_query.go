package query

import (
	"context"

	"github.com/janicaleksander/cloud/claimservice/domain"
	"github.com/mehdihadeli/go-mediatr"
)

type GetClaimsQuery struct{}

type GetClaimsQueryResponse struct {
	Claims []*GetClaimByIdQueryResponse
}

type GetClaimsQueryHandler struct {
	repo domain.ClaimRepository
}

func NewGetClaimsQueryHandler(r domain.ClaimRepository) *GetClaimsQueryHandler {
	return &GetClaimsQueryHandler{repo: r}
}
func (h *GetClaimsQueryHandler) SelfRegister() error {
	return mediatr.RegisterRequestHandler[*GetClaimsQuery, *GetClaimsQueryResponse](h)
}

func (h *GetClaimsQueryHandler) Handle(ctx context.Context, query *GetClaimsQuery) (*GetClaimsQueryResponse, error) {
	claimsDomains, err := h.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	claimsResponse := make([]*GetClaimByIdQueryResponse, len(claimsDomains))
	for i, claimDomain := range claimsDomains {
		filesResponse := make([]FileResponse, len(claimDomain.Files))
		for j, file := range claimDomain.Files {
			filesResponse[j] = FileResponse{
				ID:         file.ID,
				FileName:   file.FileName,
				FileExt:    file.FileExt,
				FileSize:   file.FileSize,
				UploadedAt: file.UploadedAt,
				StorageURL: file.StorageURL,
			}
		}
		claimsResponse[i] = &GetClaimByIdQueryResponse{
			ID:           claimDomain.ID,
			UserID:       claimDomain.UserID,
			Email:        claimDomain.Email,
			VIN:          claimDomain.VIN,
			AccidentDate: claimDomain.AccidentDate,
			Status:       string(claimDomain.Status),
			Files:        filesResponse,
			UpdatedAt:    claimDomain.UpdatedAt,
		}
	}
	return &GetClaimsQueryResponse{Claims: claimsResponse}, nil
}
