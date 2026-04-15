package meta_file_persistence

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/janicaleksander/cloud/claimservice/domain"
)

const tableName = "meta_file"

type MetaFileRepository struct {
	client *dynamodb.Client
}

func NewMetaFileRepository(client *dynamodb.Client) *MetaFileRepository {
	return &MetaFileRepository{client: client}
}

func (m *MetaFileRepository) Create(ctx context.Context, metaFile *domain.MetaFile) (*domain.MetaFile, error) {
	metaFileModel := MetaFileDomainToModel(metaFile)
	item, err := attributevalue.MarshalMap(metaFileModel)
	if err != nil {
		return nil, err
	}

	_, err = m.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
	if err != nil {
		return nil, err
	}
	return metaFile, nil
}

func (m *MetaFileRepository) GetFileById(ctx context.Context, id string) (*domain.MetaFile, error) {
	result, err := m.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, errors.New("file not found")
	}

	var model MetaFileModel
	err = attributevalue.UnmarshalMap(result.Item, &model)
	if err != nil {
		return nil, err
	}

	return MetaFileModelToDomain(&model), nil
}

func (m *MetaFileRepository) GetFiles(ctx context.Context) ([]*domain.MetaFile, error) {
	result, err := m.client.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		return nil, err
	}

	var items []MetaFileModel
	err = attributevalue.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		return nil, err
	}

	domains := make([]*domain.MetaFile, len(items))
	for i, item := range items {
		domains[i] = MetaFileModelToDomain(&item)
	}
	return domains, nil
}

func (m *MetaFileRepository) DeleteFileById(ctx context.Context, id string) error {
	_, err := m.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	return err
}
