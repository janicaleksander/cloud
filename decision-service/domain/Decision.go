package domain

type DecisionResult string

const (
	WAITING  DecisionResult = "WAITING"
	ACCEPTED DecisionResult = "ACCEPTED"
	REJECTED DecisionResult = "REJECTED"
)

type Decision struct {
	ID         uint
	ClaimID    uint
	EmployeeID *uint
	Result     DecisionResult
	Payout     float64
}

type DecisionRepository interface {
	Save(decision *Decision) (*Decision, error)
	GetAll() ([]*Decision, error)
	GetAllWaiting() ([]*Decision, error)
	GetByID(id uint) (*Decision, error)
	Update(decision *Decision) (*Decision, error)
	DeleteById(id uint) error
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
