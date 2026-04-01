package domain

type Notification struct {
}

type NotificationReceiver struct {
	ID      uint
	ClaimID uint
	Email   string
}

// todo ctx and etc
type NotificationRepository interface {
	//SaveNotification(notification *Notification) error
	//GetNotificationsByClaimID(userID uint) ([]*Notification, error)
	GetEmailByClaimID(uint) (string, error)
	SaveNotificationReceiver(receiver *NotificationReceiver) error
	UpdateNotificationReceiver(receiver *NotificationReceiver) error
	//GetNotificationReceiversByUserID(userID uint) ([]*NotificationReceiver, error)
}
