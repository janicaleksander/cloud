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

var TableNameNotification = "notification_table"
var TableNameNotificationReceiver = "notification_receiver_table"

type TableDB struct {
	Client *dynamodb.Client
}

func NewTableDB() (*TableDB, error) {
	slog.Info("initializing DynamoDB client")
	awsRegion := os.Getenv("AWS_REGION")
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(awsRegion))

	if err != nil {
		slog.Error("unable to load SDK config, %v", err)
		return nil, err
	}

	svc := dynamodb.NewFromConfig(cfg)
	_, err = svc.ListTables(context.TODO(), &dynamodb.ListTablesInput{
		Limit: aws.Int32(5),
	})
	if err != nil {
		slog.Error("unable to list tables, %v", err)
		return nil, err
	}
	slog.Info("DynamoDB client initialized successfully")

	return &TableDB{
		Client: svc,
	}, nil
}

func (t *TableDB) Migrate() error {
	paramNotificationTable := &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("notification_id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("claim_id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("notification_id"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("claim_id"),
				KeyType:       types.KeyTypeRange,
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
		TableName: aws.String(TableNameNotification),
	}

	var riue *types.ResourceInUseException
	_, err := t.Client.CreateTable(context.Background(), paramNotificationTable)
	if err != nil && !errors.As(err, &riue) {
		return err
	}

	paramNotificationReceiverTable := &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("notification_receiver_id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("claim_id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("notification_receiver_id"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("claim_id"),
				KeyType:       types.KeyTypeRange,
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
		TableName: aws.String(TableNameNotificationReceiver),
	}

	var riue2 *types.ResourceInUseException
	_, err = t.Client.CreateTable(context.Background(), paramNotificationReceiverTable)
	if err != nil && !errors.As(err, &riue2) {
		return err
	}
	return nil
}
