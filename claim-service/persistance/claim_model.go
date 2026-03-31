package persistance

import (
	"time"

	"gorm.io/gorm"
)

type ClaimModel struct {
	gorm.Model
	UserID       uint      `gorm:"not null"`
	CarID        uint      `gorm:"not null"`
	AccidentDate time.Time `gorm:"not null"`
	Status       string    `gorm:"not null"`
	Files        []FileModel
}

type FileModel struct {
	gorm.Model
	FileName     string `gorm:"not null"`
	FileExt      string `gorm:"not null"`
	StorageURL   string `gorm:"not null"`
	ClaimModelID uint   `gorm:"not null"`
}
