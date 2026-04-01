package persistance

import "gorm.io/gorm"

type NotificationModel struct {
	gorm.Model
	//todo add
}

type NotificationReceiverModel struct {
	gorm.Model
	ClaimID uint   `gorm:"not null"`
	Email   string `gorm:"not null"`
}
