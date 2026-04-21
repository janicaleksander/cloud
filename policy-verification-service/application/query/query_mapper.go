package query

import "github.com/janicaleksander/cloud/policyverificationservice/domain"

func PolicyDomainToQueryResponse(policyDomain *domain.Policy) *GetPolicyQueryResponse {
	return &GetPolicyQueryResponse{
		ID:     policyDomain.ID.String(),
		UserID: policyDomain.UserID.String(),
		VIN:    policyDomain.VIN,
		From:   policyDomain.From,
		To:     policyDomain.To,
	}
}
