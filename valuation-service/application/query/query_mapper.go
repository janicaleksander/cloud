package query

import "github.com/janicaleksander/cloud/valuationservice/domain"

func ValuationDomainToQueryResponse(valuationDomain *domain.Valuation) *GetValuationQueryResponse {
	parts := make([]PartResponse, len(valuationDomain.Parts))
	for i, part := range valuationDomain.Parts {
		parts[i] = *PartDomainToPartResponse(part)
	}
	return &GetValuationQueryResponse{
		ID:      valuationDomain.ID.String(),
		ClaimID: valuationDomain.ClaimID.String(),
		Amount:  valuationDomain.Amount,
		Parts:   parts,
	}
}

func PartDomainToPartResponse(partDomain *domain.Part) *PartResponse {
	return &PartResponse{
		ID:   partDomain.ID.String(),
		Name: partDomain.Name,
		Cost: partDomain.Cost,
	}
}
