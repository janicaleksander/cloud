package persistance

import (
	"time"

	"gorm.io/gorm"
)

type PolicyModel struct {
	gorm.Model
	UserID uint      `gorm:"not null"`
	VIN    string    `gorm:"not null"`
	From   time.Time `gorm:"not null"`
	To     time.Time `gorm:"not null"`
}
