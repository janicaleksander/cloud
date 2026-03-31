package persistance

import "gorm.io/gorm"

type ValuationModel struct {
	gorm.Model
	ClaimID uint    `gorm:"not null"`
	Amount  float64 `gorm:"not null"`
}
