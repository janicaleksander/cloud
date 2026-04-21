package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/valuationservice/domain"
	"github.com/mehdihadeli/go-mediatr"
)

type GetValuationQuery struct {
	ValuationID string
}

type GetValuationQueryResponse struct {
	ID      string         `json:"id"`
	ClaimID string         `json:"claim_id"`
	Amount  float64        `json:"amount"`
	Parts   []PartResponse `json:"parts"`
}

type PartResponse struct {
	ID   string  `json:"id"`
	Name string  `json:"name"`
	Cost float64 `json:"cost"`
}
type GetValuationQueryHandler struct {
	repo domain.ValuationRepository
}

func NewGetValuationQueryHandler(r domain.ValuationRepository) *GetValuationQueryHandler {
	return &GetValuationQueryHandler{repo: r}
}

func (h *GetValuationQueryHandler) SelfRegister() error {
	return mediatr.RegisterRequestHandler[*GetValuationQuery, *GetValuationQueryResponse](h)
}

func (h *GetValuationQueryHandler) Handle(ctx context.Context, query *GetValuationQuery) (*GetValuationQueryResponse, error) {
	cid, err := uuid.Parse(query.ValuationID)
	if err != nil {
		return nil, err
	}
	valuationDomain, err := h.repo.GetById(ctx, cid)
	if err != nil {
		return nil, err
	}
	return ValuationDomainToQueryResponse(valuationDomain), nil

}
