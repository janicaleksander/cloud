package persistence

import (
	"context"
	"errors"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/notificationservice/domain"
	"github.com/janicaleksander/cloud/notificationservice/infrastructure/tableDB"
)

type NotificationRepository struct {
	db *tableDB.TableDB
}

func NewNotificationRepository(db *tableDB.TableDB) *NotificationRepository {
	slog.Info("Initializing NotificationRepository")
	return &NotificationRepository{
		db: db,
	}
}
func (nr *NotificationRepository) SaveNotification(ctx context.Context, notification *domain.Notification) (*domain.Notification, error) {
	slog.Info("Saving notification to database", "notification", notification)
	notificationModel := NotificationDomainToModel(notification)
	av, err := attributevalue.MarshalMap(notificationModel)
	if err != nil {
		return nil, err
	}
	_, err = nr.db.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableDB.TableNameNotification),
		Item:      av,
	})
	if err != nil {
		return nil, err
	}
	return notification, nil
}

func (nr *NotificationRepository) GetNotification(ctx context.Context, id uuid.UUID) (*domain.Notification, error) {
	slog.Info("Getting notification from database", "id", id)
	input := &dynamodb.QueryInput{
		TableName: aws.String(tableDB.TableNameNotification),
		KeyConditionExpression: aws.String(
			"notification_id = :pk",
		),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: id.String()},
		},
	}
	response, err := nr.db.Client.Query(ctx, input)
	if err != nil {
		return nil, err
	}
	if len(response.Items) == 0 {
		return nil, errors.New("notification not found")
	}
	if len(response.Items) > 1 {
		return nil, errors.New("multiple notifications found with the same ID")
	}
	notificationDomain, err := NotificationModelToDomain(response.Items[0])
	return notificationDomain, err

}

func (nr *NotificationRepository) GetNotifications(ctx context.Context) ([]*domain.Notification, error) {
	slog.Info("Getting all notifications from database")
	items := make([]map[string]types.AttributeValue, 0)
	var lastKey map[string]types.AttributeValue
	for {
		response, err := nr.db.Client.Scan(ctx, &dynamodb.ScanInput{
			TableName:         aws.String(tableDB.TableNameNotification),
			ExclusiveStartKey: lastKey,
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
	notificationDomains := make([]*domain.Notification, 0)
	for _, item := range items {
		nnotification, err := NotificationModelToDomain(item)
		if err != nil {
			slog.Error("Error converting notification model to domain", "error", err)
			continue
		}
		notificationDomains = append(notificationDomains, nnotification)
	}
	return notificationDomains, nil
}

func (nr *NotificationRepository) GetNotificationsByClaimID(ctx context.Context, claimID uuid.UUID) ([]*domain.Notification, error) {
	slog.Info("Getting notifications by claim ID from database", "claimID", claimID)
	items := make([]map[string]types.AttributeValue, 0)
	var lastKey map[string]types.AttributeValue
	for {
		input := &dynamodb.QueryInput{
			TableName:              aws.String(tableDB.TableNameNotification),
			IndexName:              aws.String("user_id-index"),
			KeyConditionExpression: aws.String("claim_id = :claimID"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":claimID": &types.AttributeValueMemberS{Value: claimID.String()},
			},
			ExclusiveStartKey: lastKey,
		}
		response, err := nr.db.Client.Query(ctx, input)
		if err != nil {
			return nil, err
		}
		items = append(items, response.Items...)
		if response.LastEvaluatedKey == nil {
			break
		}
		lastKey = response.LastEvaluatedKey
	}
	notificationDomains := make([]*domain.Notification, 0)
	for _, item := range items {
		nnotification, err := NotificationModelToDomain(item)
		if err != nil {
			slog.Error("Error converting notification model to domain", "error", err)
			continue
		}
		notificationDomains = append(notificationDomains, nnotification)
	}
	return notificationDomains, nil
}
func (nr *NotificationRepository) DeleteNotificationByID(ctx context.Context, notID uuid.UUID) error {
	slog.Info("Deleting notification by ID from database", "notID", notID)
	nnotification, err := nr.GetNotification(ctx, notID)
	if err != nil {
		return err
	}
	if nnotification == nil {
		return errors.New("notification not found")
	}
	_, err = nr.db.Client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(tableDB.TableNameNotification),
		Key: map[string]types.AttributeValue{
			"notification_id": &types.AttributeValueMemberS{Value: notID.String()},
			"claim_id":        &types.AttributeValueMemberS{Value: nnotification.ClaimID.String()},
		},
	})
	return err
}

func (nr *NotificationRepository) SaveNotificationReceiver(ctx context.Context, receiver *domain.NotificationReceiver) (*domain.NotificationReceiver, error) {
	slog.Info("Saving notification receiver to database", "receiver", receiver)
	receiverModel := NotificationReceiverDomainToModel(receiver)
	av, err := attributevalue.MarshalMap(receiverModel)
	if err != nil {
		return nil, err
	}
	_, err = nr.db.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableDB.TableNameNotificationReceiver),
		Item:      av,
	})
	if err != nil {
		return nil, err
	}
	return receiver, nil
}

func (nr *NotificationRepository) GetEmailByClaimID(ctx context.Context, claimID uuid.UUID) (string, error) {
	slog.Info("Getting email by claim ID from database", "claimID", claimID)
	input := &dynamodb.QueryInput{
		TableName: aws.String(tableDB.TableNameNotificationReceiver),
		IndexName: aws.String("claim_id-index"),
		KeyConditionExpression: aws.String(
			"claim_id = :sk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":sk": &types.AttributeValueMemberS{Value: claimID.String()},
		},
	}
	response, err := nr.db.Client.Query(ctx, input)
	if err != nil {
		return "", err
	}
	if len(response.Items) == 0 {
		return "nil", errors.New("claimID not found")
	}
	if len(response.Items) > 1 {
		return "", errors.New("multiple emails found with the same claim ID")
	}
	r, err := NotificationReceiverModelToDomain(response.Items[0])
	if err != nil {
		return "", err
	}
	return r.Email, nil
}
