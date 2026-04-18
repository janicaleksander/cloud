package persistance

import (
	"github.com/janicaleksander/cloud/policyverificationservice/domain"
)

// model -> domain

func PolicyModelToDomain(pm *PolicyModel) *domain.Policy {
	return &domain.Policy{
		ID:     pm.ID,
		UserID: pm.UserID,
		VIN:    pm.VIN,
		From:   pm.From,
		To:     pm.To,
	}
}

//domain -> model

func PolicyDomainToModel(pd *domain.Policy) *PolicyModel {
	return &PolicyModel{
		ID:     pd.ID,
		UserID: pd.UserID,
		VIN:    pd.VIN,
		From:   pd.From,
		To:     pd.To,
	}
}
