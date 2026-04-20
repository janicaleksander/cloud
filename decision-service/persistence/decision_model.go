package persistence

import (
	"github.com/google/uuid"
)

type DecisionModel struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	ClaimID    uuid.UUID `gorm:"not null"`
	EmployeeID uuid.UUID `gorm:""`
	Payout     float64   `gorm:"not null"`
	Result     string    `gorm:"not null"`
}
