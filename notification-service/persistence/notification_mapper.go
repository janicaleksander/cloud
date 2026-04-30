package persistence

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/notificationservice/domain"
)

func NotificationReceiverModelToDomain(row map[string]types.AttributeValue) (*domain.NotificationReceiver, error) {
	var notificationReceiver NotificationReceiverModel
	err := attributevalue.UnmarshalMap(row, &notificationReceiver)
	if err != nil {
		return nil, err
	}
	rid, err := uuid.Parse(notificationReceiver.ID)
	if err != nil {
		return nil, err
	}
	cid, err := uuid.Parse(notificationReceiver.ClaimID)
	if err != nil {
		return nil, err
	}
	return &domain.NotificationReceiver{
		ID:      rid,
		ClaimID: cid,
		Email:   notificationReceiver.Email,
	}, nil
}

func NotificationReceiverDomainToModel(receiver *domain.NotificationReceiver) *NotificationReceiverModel {
	return &NotificationReceiverModel{
		ID:      receiver.ID.String(),
		ClaimID: receiver.ClaimID.String(),
		Email:   receiver.Email,
	}
}

func NotificationModelToDomain(row map[string]types.AttributeValue) (*domain.Notification, error) {
	var nnotification NotificationModel
	err := attributevalue.UnmarshalMap(row, &nnotification)
	if err != nil {
		return nil, err
	}

	nid, err := uuid.Parse(nnotification.ID)
	if err != nil {
		return nil, err
	}
	cid, err := uuid.Parse(nnotification.ClaimID)
	if err != nil {
		return nil, err
	}
	return &domain.Notification{
		ID:      nid,
		ClaimID: cid,
		Body:    nnotification.Body,
		SentTo:  nnotification.SentTo,
		Time:    nnotification.Time,
	}, nil
}

func NotificationDomainToModel(notification *domain.Notification) *NotificationModel {
	return &NotificationModel{
		ID:      notification.ID.String(),
		ClaimID: notification.ClaimID.String(),
		Body:    notification.Body,
		SentTo:  notification.SentTo,
		Time:    notification.Time,
	}
}
