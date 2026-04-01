package application

import "github.com/janicaleksander/cloud/notificationservice/domain"

type NotificationService struct {
	notificationRepository domain.NotificationRepository
}

func NewNotificationService(notificationRepository domain.NotificationRepository) *NotificationService {
	return &NotificationService{
		notificationRepository: notificationRepository,
	}
}

func (s *NotificationService) CreateNotificationReceiver(nr *domain.NotificationReceiver) error {
	return s.notificationRepository.SaveNotificationReceiver(nr)
}

func (s *NotificationService) UpdateNotificationReceiver(nr *domain.NotificationReceiver) error {
	return s.notificationRepository.UpdateNotificationReceiver(nr)
}

func (s *NotificationService) GetEmailByClaimID(claimID uint) (string, error) {
	email, err := s.notificationRepository.GetEmailByClaimID(claimID)
	if err != nil {
		return "", err
	}
	return email, nil
}
