package persistance

import (
	"gorm.io/gorm"
)

type DecisionModel struct {
	gorm.Model
	ClaimID    uint    `gorm:"not null"`
	EmployeeID *uint   `gorm:""`
	Payout     float64 `gorm:"not null"`
	Result     string  `gorm:"not null"`
}
