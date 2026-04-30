package persistance

import (
	"context"
	"errors"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/policyverificationservice/domain"
	"github.com/janicaleksander/cloud/policyverificationservice/infrastructure/tableDB"
)

type PolicyRepository struct {
	db *tableDB.TableDB
}

func NewPolicyRepository(tableDB *tableDB.TableDB) *PolicyRepository {
	slog.Info("Initializing PolicyRepository")
	return &PolicyRepository{db: tableDB}
}

func (pr *PolicyRepository) GetAll(ctx context.Context) ([]*domain.Policy, error) {
	slog.Info("Getting all policies from the database")
	items := make([]map[string]types.AttributeValue, 0, 32)
	var lastKey map[string]types.AttributeValue

	for {
		response, err := pr.db.Client.Scan(
			ctx, &dynamodb.ScanInput{
				ExclusiveStartKey: lastKey,
				TableName:         aws.String(tableDB.TableNamePolicy),
			})
		if err != nil {
			return nil, err
		}
		items = append(items, response.Items...)
		if response.LastEvaluatedKey == nil {
			break
		}
		lastKey = response.LastEvaluatedKey
	}

	domainPolicies := make([]*domain.Policy, 0, len(items))
	for idx := range items {
		policy, err := PolicyModelToDomain(items[idx])
		if err != nil {
			slog.Error("Error converting policy model to domain", "error", err)
			continue
		}
		domainPolicies = append(domainPolicies, policy)
	}
	return domainPolicies, nil

}

func (pr *PolicyRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.Policy, error) {
	slog.Info("Getting policy by ID from the database", "id", id)
	input := &dynamodb.QueryInput{
		TableName: aws.String(tableDB.TableNamePolicy),
		KeyConditionExpression: aws.String(
			"policy_id = :pk",
		),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: id.String()},
		},
	}
	response, err := pr.db.Client.Query(ctx, input)
	if err != nil {
		return nil, err
	}
	if len(response.Items) == 0 {
		return nil, errors.New("policy not found")
	}
	if len(response.Items) > 1 {
		slog.Warn("Multiple policies found with the same ID, returning the first one", "id", id)
		return nil, errors.New("multiple policies found with the same ID")
	}
	policyDomain, err := PolicyModelToDomain(response.Items[0])
	return policyDomain, err

}

func (pr *PolicyRepository) Save(ctx context.Context, p *domain.Policy) (*domain.Policy, error) {
	slog.Info("Saving policy to the database")
	policyModel := PolicyDomainToModel(p)
	av, err := attributevalue.MarshalMap(policyModel)
	if err != nil {
		return nil, err
	}
	_, err = pr.db.Client.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableDB.TableNamePolicy),
	})
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (pr *PolicyRepository) Update(ctx context.Context, p *domain.Policy) (*domain.Policy, error) {
	slog.Info("Updating policy in the database", "policy", p)
	policyModel := PolicyDomainToModel(p)

	av, err := attributevalue.MarshalMap(policyModel)
	if err != nil {
		return nil, err
	}

	_, err = pr.db.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(tableDB.TableNamePolicy),
		Item:                av,
		ConditionExpression: aws.String("attribute_exists(policy_id)"),
	})
	if err != nil {
		return nil, err
	}
	return p, nil

}
func (pr *PolicyRepository) DeleteById(ctx context.Context, id uuid.UUID) error {
	slog.Info("Deleting policy by ID from the database", "id", id)
	policy, err := pr.GetById(ctx, id)
	if err != nil {
		slog.Error("Error getting policy by ID before deletion", "id", id, "error", err)
		return err
	}
	if policy == nil {
		return errors.New("policy not found")
	}
	_, err = pr.db.Client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(tableDB.TableNamePolicy),
		Key: map[string]types.AttributeValue{
			"policy_id": &types.AttributeValueMemberS{Value: id.String()},
			"user_id":   &types.AttributeValueMemberS{Value: policy.UserID.String()},
		},
	})
	return err

}
func (pr *PolicyRepository) IfUserHasPolicy(ctx context.Context, userID uuid.UUID, vin string) (bool, *domain.Policy) {
	slog.Info("Checking if user has policy for given VIN", "userID", userID, "vin", vin)
	response, err := pr.db.Client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(tableDB.TableNamePolicy),
		IndexName:              aws.String("user_id-index"),
		KeyConditionExpression: aws.String("user_id = :userID"),
		FilterExpression:       aws.String("vin = :Vin"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userID": &types.AttributeValueMemberS{Value: userID.String()},
			":Vin":    &types.AttributeValueMemberS{Value: vin},
		},
	})
	if err != nil {
		slog.Error("Error querying policies by user ID and VIN", "error", err)
		return false, nil
	}
	if len(response.Items) == 0 {
		return false, nil
	}
	policy, err := PolicyModelToDomain(response.Items[0])
	if err != nil {
		slog.Error("Error converting policy model to domain", "error", err)
		return false, nil
	}
	return true, policy
}
