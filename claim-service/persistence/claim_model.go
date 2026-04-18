package persistence

import (
	"time"

	"github.com/google/uuid"
)

type ClaimModel struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID       uuid.UUID `gorm:"not null"`
	Email        string    `gorm:"not null"`
	VIN          string    `gorm:"not null"`
	AccidentDate time.Time `gorm:"not null"`
	Status       string    `gorm:"not null"`
	Files        []FileModel

	UpdatedAt time.Time
}

type FileModel struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	FileName     string    `gorm:"not null"`
	FileExt      string    `gorm:"not null"`
	FileSize     int64     `gorm:"not null"`
	StorageURL   string    `gorm:"not null"`
	ClaimModelID uuid.UUID `gorm:"not null"`
	UploadedAt   time.Time `gorm:"not null"`
}
