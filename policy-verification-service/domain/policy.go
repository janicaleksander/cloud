package domain

import (
	"context"
	"time"
)

type RejectionReason string

const (
	InvalidVIN               RejectionReason = "Invalid VIN"
	PolicyNotFound           RejectionReason = "Policy not found"
	PolicyExpired            RejectionReason = "Policy expired"
	AccidentDateBeforePolicy RejectionReason = "Accident date is before policy start date"
)

type Policy struct {
	ID        uint
	UserID    uint
	VIN       string
	From      time.Time
	To        time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
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

//TODO This repo

type PolicyRepository interface {
	GetAll(context.Context) ([]*Policy, error)
	GetById(context.Context, uint) (*Policy, error)
	Save(context.Context, *Policy) (*Policy, error)
	Update(context.Context, *Policy) (*Policy, error)
	DeleteById(context.Context, uint) error
	IfUserHasPolicy(context.Context, uint, string) (bool, *Policy)
}
