package persistance

import (
	"context"
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
	tableDB *tableDB.TableDB
}

func NewPolicyRepository(tableDB *tableDB.TableDB) *PolicyRepository {
	slog.Info("Initializing PolicyRepository")
	return &PolicyRepository{tableDB: tableDB}
}

func (pr *PolicyRepository) GetAll(ctx context.Context) ([]*domain.Policy, error) {
	slog.Info("Getting all policies from the database")
	items := make([]map[string]types.AttributeValue, 0, 32)
	var lastKey map[string]types.AttributeValue

	for {
		response, err := pr.tableDB.Client.Scan(
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
	response, err := pr.tableDB.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableDB.TableNamePolicy),
		Key: map[string]types.AttributeValue{
			"policy_id": &types.AttributeValueMemberS{Value: id.String()},
		},
	})
	if err != nil {
		return nil, err
	}
	policyDomain, err := PolicyModelToDomain(response.Item)
	return policyDomain, err

}

func (pr *PolicyRepository) Save(ctx context.Context, p *domain.Policy) (*domain.Policy, error) {
	slog.Info("Saving policy to the database")
	policyModel := PolicyDomainToModel(p)
	av, err := attributevalue.MarshalMap(policyModel)
	if err != nil {
		return nil, err
	}
	_, err = pr.tableDB.Client.PutItem(ctx, &dynamodb.PutItemInput{
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

	_, err = pr.tableDB.Client.PutItem(ctx, &dynamodb.PutItemInput{
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
	_, err := pr.tableDB.Client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(tableDB.TableNamePolicy),
		Key: map[string]types.AttributeValue{
			"policy_id": &types.AttributeValueMemberS{Value: id.String()},
		},
	})
	return err

}
func (pr *PolicyRepository) IfUserHasPolicy(ctx context.Context, userID uuid.UUID, vin string) (bool, *domain.Policy) {
	slog.Info("Checking if user has policy for given VIN", "userID", userID, "vin", vin)
	response, err := pr.tableDB.Client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(tableDB.TableNamePolicy),
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
