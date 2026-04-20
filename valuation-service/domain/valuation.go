package domain

import (
	"context"

	"github.com/google/uuid"
)

type Valuation struct {
	ID      uuid.UUID
	ClaimID uuid.UUID
	Amount  float64
	Parts   []*Part
}

type Part struct {
	ID   uuid.UUID
	Name string
	Cost float64
}
type ValuationRepository interface {
	GetAll(context.Context) ([]*Valuation, error)
	GetById(context.Context, uuid.UUID) (*Valuation, error)
	Save(context.Context, *Valuation) (*Valuation, error)
	Update(context.Context, *Valuation) (*Valuation, error)
	DeleteById(context.Context, uuid.UUID) error
}
