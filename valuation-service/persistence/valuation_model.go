package persistence

import (
	"github.com/google/uuid"
)

type ValuationModel struct {
	ID      uuid.UUID   `gorm:"type:uuid;primaryKey"`
	ClaimID uuid.UUID   `gorm:"not null;index"`
	Amount  float64     `gorm:"not null"`
	Parts   []PartModel `gorm:"foreignKey:ValuationID"`
}

type PartModel struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	ValuationID uuid.UUID `gorm:"not null;index"`
	Name        string    `gorm:"not null;size:255"`
	Cost        float64   `gorm:"not null"`
}
