package domain

import (
	"context"
	"time"
)

type Valuation struct {
	ID        uint
	ClaimID   uint
	Amount    float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ValuationRepository interface {
	GetAll(context.Context) ([]*Valuation, error)
	GetById(context.Context, uint) (*Valuation, error)
	Save(context.Context, *Valuation) (*Valuation, error)
	Update(context.Context, *Valuation) (*Valuation, error)
	DeleteById(context.Context, uint) error
}
