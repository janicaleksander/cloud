package persistance

import (
	"context"

	"github.com/janicaleksander/cloud/notificationservice/domain"
	"gorm.io/gorm"
)

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{
		db: db,
	}
}
func (nr *NotificationRepository) SaveNotification(ctx context.Context, notification *domain.Notification) (*domain.Notification, error) {
	notificationModel := NotificationDomainToModel(notification)
	err := gorm.G[NotificationModel](nr.db).Create(ctx, notificationModel)
	if err != nil {
		return nil, err
	}
	return NotificationModelToDomain(notificationModel), nil

}

func (nr *NotificationRepository) GetNotification(ctx context.Context, id uint) (*domain.Notification, error) {
	notificationModel, err := gorm.G[NotificationModel](nr.db).Where("id = ? ", id).First(ctx)
	if err != nil {
		return nil, err
	}
	return NotificationModelToDomain(&notificationModel), nil
}

func (nr *NotificationRepository) GetNotifications(ctx context.Context) ([]*domain.Notification, error) {
	notificationDomains := make([]*domain.Notification, 0)
	notificationModels, err := gorm.G[NotificationModel](nr.db).Find(ctx)
	if err != nil {
		return nil, err
	}
	for _, model := range notificationModels {
		notificationDomains = append(notificationDomains, NotificationModelToDomain(&model))
	}
	return notificationDomains, nil
}

func (nr *NotificationRepository) GetNotificationsByClaimID(ctx context.Context, claimID uint) ([]*domain.Notification, error) {
	notificationDomains := make([]*domain.Notification, 0)
	notificationModels, err := gorm.G[NotificationModel](nr.db).Where("claim_id = ?", claimID).Find(ctx)
	if err != nil {
		return nil, err
	}
	for _, model := range notificationModels {
		notificationDomains = append(notificationDomains, NotificationModelToDomain(&model))
	}
	return notificationDomains, nil

}
func (nr *NotificationRepository) DeleteNotificationByID(ctx context.Context, notID uint) error {
	_, err := gorm.G[NotificationModel](nr.db).Where("id = ?", notID).Delete(ctx)
	return err
}

func (nr *NotificationRepository) SaveNotificationReceiver(ctx context.Context, receiver *domain.NotificationReceiver) (*domain.NotificationReceiver, error) {
	notificationReceiverModel := NotificationReceiverDomainToModel(receiver)
	err := gorm.G[NotificationReceiverModel](nr.db).Create(ctx, notificationReceiverModel)
	if err != nil {
		return nil, err
	}
	return NotificationReceiverModelToDomain(notificationReceiverModel), nil
}

func (nr *NotificationRepository) UpdateNotificationReceiver(ctx context.Context, receiver *domain.NotificationReceiver) (*domain.NotificationReceiver, error) {
	model := NotificationReceiverDomainToModel(receiver)
	err := nr.db.Save(model).Error
	if err != nil {
		return nil, err
	}
	return NotificationReceiverModelToDomain(model), nil
}

func (nr *NotificationRepository) GetEmailByClaimID(ctx context.Context, claimID uint) (string, error) {
	var receiver NotificationReceiverModel
	err := nr.db.Where("claim_id = ?", claimID).First(&receiver).Error
	if err != nil {
		return "", err
	}
	return receiver.Email, nil
}
