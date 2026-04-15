package claim_persistence

import (
	"time"

	"gorm.io/gorm"
)

type ClaimModel struct {
	gorm.Model
	UserID       uint      `gorm:"not null"`
	Email        string    `gorm:"not null"`
	VIN          string    `gorm:"not null"`
	AccidentDate time.Time `gorm:"not null"`
	Status       string    `gorm:"not null"`
	Files        []FileModel
}

type FileModel struct {
	gorm.Model
	StorageURL   string `gorm:"not null"`
	ClaimModelID uint   `gorm:"not null"`
}
