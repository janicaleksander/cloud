package persistance

import "gorm.io/gorm"

type ValuationModel struct {
	gorm.Model
	ClaimID uint        `gorm:"not null"`
	Amount  float64     `gorm:"not null"`
	Parts   []PartModel `gorm:"many2many:valuation_parts;"`
}

type PartModel struct {
	gorm.Model
	Name string  `gorm:"not null"`
	Cost float64 `gorm:"not null"`
}
