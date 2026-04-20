package domain

import (
	"context"

	"github.com/google/uuid"
)

type DecisionResult string

const (
	WAITING  DecisionResult = "WAITING"
	ACCEPTED DecisionResult = "ACCEPTED"
	REJECTED DecisionResult = "REJECTED"
)

type Decision struct {
	ID         uuid.UUID
	ClaimID    uuid.UUID
	EmployeeID uuid.UUID
	Result     DecisionResult
	Payout     float64
}

type DecisionRepository interface {
	Save(ctx context.Context, decision *Decision) (*Decision, error)
	GetAll(ctx context.Context) ([]*Decision, error)
	GetAllWaiting(ctx context.Context) ([]*Decision, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Decision, error)
	Update(ctx context.Context, decision *Decision) (*Decision, error)
	DeleteById(ctx context.Context, id uuid.UUID) error
}

func StringToResult(s string) DecisionResult {
	switch s {
	case string(WAITING):
		return WAITING
	case string(ACCEPTED):
		return ACCEPTED
	case string(REJECTED):
		return REJECTED
	default:
		return ""
	}
}
