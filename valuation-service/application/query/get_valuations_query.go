package query

import (
	"context"

	"github.com/janicaleksander/cloud/valuationservice/domain"
	"github.com/mehdihadeli/go-mediatr"
)

type GetValuationsQuery struct{}

type GetValuationsQueryResponse struct {
	Valuations []*GetValuationQueryResponse
}

type GetValuationsQueryHandler struct {
	repo domain.ValuationRepository
}

func NewGetValuationsQueryHandler(r domain.ValuationRepository) *GetValuationsQueryHandler {
	return &GetValuationsQueryHandler{repo: r}
}

func (h *GetValuationsQueryHandler) SelfRegister() error {
	return mediatr.RegisterRequestHandler[*GetValuationsQuery, *GetValuationsQueryResponse](h)
}

func (h *GetValuationsQueryHandler) Handle(ctx context.Context, query *GetValuationsQuery) (*GetValuationsQueryResponse, error) {
	valuationDomains, err := h.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	responses := make([]*GetValuationQueryResponse, len(valuationDomains))
	for i, valuationDomain := range valuationDomains {
		parts := make([]PartResponse, len(valuationDomain.Parts))
		for j, part := range valuationDomain.Parts {
			parts[j] = PartResponse{
				ID:   part.ID.String(),
				Name: part.Name,
				Cost: part.Cost,
			}
		}
		responses[i] = &GetValuationQueryResponse{
			ID:      valuationDomain.ID.String(),
			ClaimID: valuationDomain.ClaimID.String(),
			Amount:  valuationDomain.Amount,
			Parts:   parts,
		}
	}
	return &GetValuationsQueryResponse{Valuations: responses}, nil
}
