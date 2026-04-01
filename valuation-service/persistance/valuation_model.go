package persistance

import "gorm.io/gorm"

type ValuationModel struct {
	gorm.Model
	ClaimID uint        `gorm:"not null;index"`
	Amount  float64     `gorm:"not null"`
	Parts   []PartModel `gorm:"foreignKey:ValuationID"`
}

type PartModel struct {
	gorm.Model
	ValuationID uint    `gorm:"not null;index"`
	Name        string  `gorm:"not null;size:255"`
	Cost        float64 `gorm:"not null"`
}
