package persistance

import (
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

func (nr *NotificationRepository) SaveNotificationReceiver(receiver *domain.NotificationReceiver) error {
	return nr.db.Create(NotificationReceiverDomainToModel(receiver)).Error
}

func (nr *NotificationRepository) UpdateNotificationReceiver(receiver *domain.NotificationReceiver) error {
	model := NotificationReceiverDomainToModel(receiver)
	return nr.db.Save(model).Error
}

func (nr *NotificationRepository) GetEmailByClaimID(claimID uint) (string, error) {
	var receiver NotificationReceiverModel
	err := nr.db.Where("claim_id = ?", claimID).First(&receiver).Error
	if err != nil {
		return "", err
	}
	return receiver.Email, nil
}
