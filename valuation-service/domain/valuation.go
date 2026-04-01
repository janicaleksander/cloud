package domain

import (
	"context"
)

type Valuation struct {
	ID      uint
	ClaimID uint
	Amount  float64
	Parts   []*Part
}

type Part struct {
	ID   uint
	Name string
	Cost float64
}
type ValuationRepository interface {
	GetAll(context.Context) ([]*Valuation, error)
	GetById(context.Context, uint) (*Valuation, error)
	Save(context.Context, *Valuation) (*Valuation, error)
	Update(context.Context, *Valuation) (*Valuation, error)
	DeleteById(context.Context, uint) error
}

//tod oprealod parts
