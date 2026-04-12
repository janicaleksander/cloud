package application

import (
	"context"
	"log/slog"

	"github.com/janicaleksander/cloud/notificationservice/domain"
)

type NotificationService struct {
	notificationRepository domain.NotificationRepository
}

func NewNotificationService(notificationRepository domain.NotificationRepository) *NotificationService {
	slog.Info("Creating NotificationService")
	return &NotificationService{
		notificationRepository: notificationRepository,
	}
}

func (s *NotificationService) CreateNotification(n *domain.Notification) (*domain.Notification, error) {
	slog.Info("Creating notification", "claimID", n.ClaimID)
	return s.notificationRepository.SaveNotification(context.Background(), n)
}

func (s *NotificationService) GetNotifications() ([]*domain.Notification, error) {
	slog.Info("Getting all notifications")
	return s.notificationRepository.GetNotifications(context.Background())
}

func (s *NotificationService) GetNotification(id uint) (*domain.Notification, error) {
	slog.Info("Getting notification with ID", "id", id)
	return s.notificationRepository.GetNotification(context.Background(), id)
}
func (s *NotificationService) GetNotificationsForClaimID(claimID uint) ([]*domain.Notification, error) {
	slog.Info("Getting notifications for claim ID", "claimID", claimID)
	return s.notificationRepository.GetNotificationsByClaimID(context.Background(), claimID)
}
func (s *NotificationService) DeleteNotification(id uint) error {
	slog.Info("Deleting notification with ID", "id", id)
	return s.notificationRepository.DeleteNotificationByID(context.Background(), id)
}

func (s *NotificationService) CreateNotificationReceiver(nr *domain.NotificationReceiver) (*domain.NotificationReceiver, error) {
	slog.Info("Creating notification receiver", "claimID", nr.ClaimID, "email", nr.Email)
	return s.notificationRepository.SaveNotificationReceiver(context.Background(), nr)
}

func (s *NotificationService) UpdateNotificationReceiver(nr *domain.NotificationReceiver) (*domain.NotificationReceiver, error) {
	slog.Info("Updating notification receiver", "claimID", nr.ClaimID, "email", nr.Email)
	return s.notificationRepository.UpdateNotificationReceiver(context.Background(), nr)
}

func (s *NotificationService) GetEmailByClaimID(claimID uint) (string, error) {
	slog.Info("Getting email by claim ID", "claimID", claimID)
	email, err := s.notificationRepository.GetEmailByClaimID(context.Background(), claimID)
	if err != nil {
		return "", err
	}
	return email, nil
}
