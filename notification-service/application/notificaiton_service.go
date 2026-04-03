package application

import (
	"context"

	"github.com/janicaleksander/cloud/notificationservice/domain"
)

type NotificationService struct {
	notificationRepository domain.NotificationRepository
}

func NewNotificationService(notificationRepository domain.NotificationRepository) *NotificationService {
	return &NotificationService{
		notificationRepository: notificationRepository,
	}
}

func (s *NotificationService) CreateNotification(n *domain.Notification) (*domain.Notification, error) {
	return s.notificationRepository.SaveNotification(context.Background(), n)
}

func (s *NotificationService) GetNotifications() ([]*domain.Notification, error) {
	return s.notificationRepository.GetNotifications(context.Background())
}

func (s *NotificationService) GetNotification(id uint) (*domain.Notification, error) {
	return s.notificationRepository.GetNotification(context.Background(), id)
}
func (s *NotificationService) GetNotificationsForClaimID(claimID uint) ([]*domain.Notification, error) {
	return s.notificationRepository.GetNotificationsByClaimID(context.Background(), claimID)
}
func (s *NotificationService) DeleteNotification(id uint) error {
	return s.notificationRepository.DeleteNotificationByID(context.Background(), id)
}

func (s *NotificationService) CreateNotificationReceiver(nr *domain.NotificationReceiver) (*domain.NotificationReceiver, error) {
	return s.notificationRepository.SaveNotificationReceiver(context.Background(), nr)
}

func (s *NotificationService) UpdateNotificationReceiver(nr *domain.NotificationReceiver) (*domain.NotificationReceiver, error) {
	return s.notificationRepository.UpdateNotificationReceiver(context.Background(), nr)
}

func (s *NotificationService) GetEmailByClaimID(claimID uint) (string, error) {
	email, err := s.notificationRepository.GetEmailByClaimID(context.Background(), claimID)
	if err != nil {
		return "", err
	}
	return email, nil
}
