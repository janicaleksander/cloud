package persistence

import (
	"time"

	"github.com/google/uuid"
)

type NotificationModel struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey"`
	ClaimID uuid.UUID `gorm:"not null"`
	Body    string    `gorm:"not null"`
	SentTo  string    `gorm:"not null"`
	Time    time.Time `gorm:"not null"`
}

type NotificationReceiverModel struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey"`
	ClaimID uuid.UUID `gorm:"not null"`
	Email   string    `gorm:"not null"`
}
