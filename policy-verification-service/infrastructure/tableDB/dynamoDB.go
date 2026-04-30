package tableDB

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var TableNamePolicy = "policy_table"

type TableDB struct {
	Client *dynamodb.Client
}

func NewTableDB() (*TableDB, error) {
	slog.Info("initializing DynamoDB client")
	awsRegion := os.Getenv("AWS_REGION")
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(awsRegion),
	)
	if err != nil {
		slog.Error("unable to load SDK config, %v", err)
		return nil, err
	}

	svc := dynamodb.NewFromConfig(cfg)
	_, err = svc.ListTables(context.TODO(), &dynamodb.ListTablesInput{
		Limit: aws.Int32(5),
	})
	if err != nil {
		return nil, err
	}
	slog.Info("DynamoDB client initialized successfully")

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
			{
				AttributeName: aws.String("user_id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("policy_id"),
				KeyType:       types.KeyTypeHash,
			}, {
				AttributeName: aws.String("user_id"),
				KeyType:       types.KeyTypeRange,
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
		TableName: aws.String(TableNamePolicy),
	}

	var riue1 *types.ResourceInUseException
	_, err := t.Client.CreateTable(context.Background(), paramPolicyTable)
	if err != nil && !errors.As(err, &riue1) {
		return err
	}
	return nil
}
