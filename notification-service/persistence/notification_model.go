package persistence

import (
	"time"

	"gorm.io/gorm"
)

type NotificationModel struct {
	gorm.Model
	ClaimID uint      `gorm:"not null"`
	Body    string    `gorm:"not null"`
	SentTo  string    `gorm:"not null"`
	Time    time.Time `gorm:"not null"`
}

type NotificationReceiverModel struct {
	gorm.Model
	ClaimID uint   `gorm:"not null"`
	Email   string `gorm:"not null"`
}
