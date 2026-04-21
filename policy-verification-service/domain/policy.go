package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type RejectionReason string

const (
	InvalidVIN               RejectionReason = "Invalid VIN"
	PolicyNotFound           RejectionReason = "Policy not found"
	PolicyExpired            RejectionReason = "Policy expired"
	AccidentDateBeforePolicy RejectionReason = "Accident date is before policy start date"
)

type Policy struct {
	ID     uuid.UUID
	UserID uuid.UUID
	VIN    string
	From   time.Time
	To     time.Time
}

func NewPolicy(id, userID uuid.UUID, vin string, from, to time.Time) *Policy {
	return &Policy{
		ID:     id,
		UserID: userID,
		VIN:    vin,
		From:   from,
		To:     to,
	}
}
func (p Policy) IsValid(accidentDate time.Time) (bool, RejectionReason) {
	if accidentDate.Before(p.From) {
		return false, AccidentDateBeforePolicy
	}
	if accidentDate.After(p.To) {
		return false, PolicyExpired
	}
	return true, ""
}

type PolicyRepository interface {
	GetAll(context.Context) ([]*Policy, error)
	GetById(context.Context, uuid.UUID) (*Policy, error)
	Save(context.Context, *Policy) (*Policy, error)
	Update(context.Context, *Policy) (*Policy, error)
	DeleteById(context.Context, uuid.UUID) error
	IfUserHasPolicy(context.Context, uuid.UUID, string) (bool, *Policy)
}
