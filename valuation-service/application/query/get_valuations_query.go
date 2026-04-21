package query

import (
	"context"

	"github.com/janicaleksander/cloud/valuationservice/domain"
	"github.com/mehdihadeli/go-mediatr"
)

type GetValuationsQuery struct{}

type GetValuationsQueryResponse struct {
	Valuations []*GetValuationQueryResponse `json:"valuations"`
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
		responses[i] = ValuationDomainToQueryResponse(valuationDomain)
	}

	return &GetValuationsQueryResponse{Valuations: responses}, nil
}
