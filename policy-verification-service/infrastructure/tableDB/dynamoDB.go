package tableDB

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var TableNamePolicy = "policy_table"
var TableNameNotification = "notification_table"

type TableDB struct {
	Client *dynamodb.Client
}

func NewTableDB() (*TableDB, error) {
	awsRegion := os.Getenv("AWS_REGION")
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(awsRegion),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	svc := dynamodb.NewFromConfig(cfg)
	_, err = svc.ListTables(context.TODO(), &dynamodb.ListTablesInput{
		Limit: aws.Int32(5),
	})
	if err != nil {
		return nil, err
	}
	return &TableDB{
		Client: svc,
	}, nil
}

func (t *TableDB) Migrate() error {
	// policy table
	paramPolicyTable := &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("policy_id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("policy_id"),
				KeyType:       types.KeyTypeHash,
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
		TableName: aws.String(TableNamePolicy),
	}
	paramNotificationTable := &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("claim_id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("default_sk"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("claim_id"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("default_sk"),
				KeyType:       types.KeyTypeRange,
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
		TableName: aws.String(TableNameNotification),
	}
	var riue1 *types.ResourceInUseException
	_, err := t.Client.CreateTable(context.Background(), paramPolicyTable)
	if err != nil && !errors.As(err, &riue1) {
		return err
	}

	var riue2 *types.ResourceInUseException
	_, err = t.Client.CreateTable(context.Background(), paramNotificationTable)
	if err != nil && !errors.As(err, &riue2) {
		return err
	}
	return nil
}
