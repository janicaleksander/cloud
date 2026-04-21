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

func NewValuation(id, claimID uuid.UUID, amount float64, parts []*Part) *Valuation {
	return &Valuation{
		ID:      id,
		ClaimID: claimID,
		Amount:  amount,
		Parts:   parts,
	}
}

func NewPart(id uuid.UUID, name string, cost float64) *Part {
	return &Part{
		ID:   id,
		Name: name,
		Cost: cost,
	}
}

type ValuationRepository interface {
	GetAll(context.Context) ([]*Valuation, error)
	GetById(context.Context, uuid.UUID) (*Valuation, error)
	Save(context.Context, *Valuation) (*Valuation, error)
	Update(context.Context, *Valuation) (*Valuation, error)
	DeleteById(context.Context, uuid.UUID) error
}
