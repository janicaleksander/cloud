package persistance

import (
	"time"

	"github.com/google/uuid"
)

type PolicyModel struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID uuid.UUID `gorm:"not null"`
	VIN    string    `gorm:"not null"`
	From   time.Time `gorm:"not null"`
	To     time.Time `gorm:"not null"`
}
