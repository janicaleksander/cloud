package persistance

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/policyverificationservice/domain"
)

// model -> domain

//domain -> model

func PolicyDomainToModel(pd *domain.Policy) *PolicyModel {
	return &PolicyModel{
		ID:     pd.ID.String(),
		UserID: pd.UserID.String(),
		VIN:    pd.VIN,
		From:   pd.From,
		To:     pd.To,
	}
}

func PolicyModelToDomain(row map[string]types.AttributeValue) (*domain.Policy, error) {
	var policy PolicyModel
	err := attributevalue.UnmarshalMap(row, &policy)
	if err != nil {
		return nil, err
	}
	pid, err := uuid.Parse(policy.ID)
	if err != nil {
		return nil, err
	}
	uid, err := uuid.Parse(policy.UserID)
	if err != nil {
		return nil, err
	}
	return &domain.Policy{
		ID:     pid,
		UserID: uid,
		VIN:    policy.VIN,
		From:   policy.From,
		To:     policy.To,
	}, nil

}
