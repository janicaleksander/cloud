package persistance

import (
	"github.com/janicaleksander/cloud/policyverificationservice/domain"
	"gorm.io/gorm"
)

// model -> domain

func PolicyModelToDomain(pm *PolicyModel) *domain.Policy {
	return &domain.Policy{
		ID:        pm.ID,
		UserID:    pm.UserID,
		VIN:       pm.VIN,
		From:      pm.From,
		To:        pm.To,
		CreatedAt: pm.CreatedAt,
		UpdatedAt: pm.UpdatedAt,
	}
}

//domain -> model

func PolicyDomainToModel(pd *domain.Policy) *PolicyModel {
	return &PolicyModel{
		Model: gorm.Model{
			ID:        pd.ID,
			CreatedAt: pd.CreatedAt,
			UpdatedAt: pd.UpdatedAt,
		},
		UserID: pd.UserID,
		VIN:    pd.VIN,
		From:   pd.From,
		To:     pd.To,
	}
}
