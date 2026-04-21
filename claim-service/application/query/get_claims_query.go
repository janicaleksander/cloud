package query

import (
	"context"

	"github.com/janicaleksander/cloud/claimservice/domain"
	"github.com/mehdihadeli/go-mediatr"
)

type GetClaimsQuery struct{}

type GetClaimsQueryResponse struct {
	Claims []*GetClaimByIdQueryResponse `json:"claims"`
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
			filesResponse[j] = *FileDomainToQueryResponse(file)
		}
		claimsResponse[i] = ClaimDomainToQueryResponse(claimDomain)
	}
	return &GetClaimsQueryResponse{Claims: claimsResponse}, nil
}
