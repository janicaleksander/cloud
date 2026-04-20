package application

/*package application

import (
	"errors"
	"log/slog"

	"github.com/janicaleksander/cloud/common/event"
	"github.com/janicaleksander/cloud/decisionservice/domain"
)

type DecisionService struct {
	decisionRepository domain.DecisionRepository
	publisher          DecisionPublisher
}

func NewDecisionService(decisionRepo domain.DecisionRepository, publisher DecisionPublisher) *DecisionService {
	slog.Info("Creating DecisionService")
	return &DecisionService{
		decisionRepository: decisionRepo,
		publisher:          publisher,
	}
}

func (ds *DecisionService) makeDecision(newDecision *domain.Decision, reason string) {
	slog.Info("Making decision for claim", "claimID", newDecision.ClaimID, "result", newDecision.Result, "employeeID", *newDecision.EmployeeID)
	if newDecision.Result == domain.ACCEPTED {
		err := ds.publisher.Publish("events", event.PayoutApprovedEvent{
			ClaimID:              newDecision.ClaimID,
			AcceptedPayoutAmount: newDecision.Payout,
			ByEmployeeID:         *newDecision.EmployeeID,
		})
		if err != nil {
			slog.Error("Failed to publish PayoutApprovedEvent", "error", err)
			return
		}
	}
	if newDecision.Result == domain.REJECTED {
		err := ds.publisher.Publish("events", event.PayoutRejectedEvent{
			ClaimID:      newDecision.ClaimID,
			Reason:       reason,
			ByEmployeeID: *newDecision.EmployeeID,
		})
		if err != nil {
			slog.Error("Failed to publish PayoutRejectedEvent", "error", err)
			return
		}
	}
}
func (ds *DecisionService) PrepareDecision(claimID uint, payoutAmount float64) (*domain.Decision, error) {
	slog.Info("Preparing decision for claim", "claimID", claimID, "payoutAmount", payoutAmount)
	decisionDomain := &domain.Decision{
		ClaimID: claimID,
		Payout:  payoutAmount,
		Result:  domain.WAITING,
	}
	createdDecision, err := ds.decisionRepository.Save(decisionDomain)
	if err != nil {
		return nil, err
	}
	return createdDecision, nil
}

func (ds *DecisionService) GetDecision(decisionID uint) (*domain.Decision, error) {
	slog.Info("Getting decision with ID", "decisionID", decisionID)
	return ds.decisionRepository.GetByID(decisionID)
}
func (ds *DecisionService) GetDecisions() ([]*domain.Decision, error) {
	slog.Info("Getting all decisions")
	return ds.decisionRepository.GetAll()

}
func (ds *DecisionService) GetWaitingDecisions() ([]*domain.Decision, error) {
	slog.Info("Getting all waiting decisions")
	return ds.decisionRepository.GetAllWaiting()
}
func (ds *DecisionService) DeleteDecision(decisionID uint) error {
	slog.Info("Deleting decision with ID", "decisionID", decisionID)
	return ds.decisionRepository.DeleteById(decisionID)

}
func (ds *DecisionService) UpdateDecisionState(oldDecision *domain.Decision, newState domain.DecisionResult, empID uint, reason string) (*domain.Decision, error) {
	slog.Info("Update Decision with ID", "decisionID", oldDecision.ID, "newState", newState, "employeeID", empID, "reason", reason)
	if oldDecision.Result != domain.WAITING {
		return nil, errors.New("already accepted/denied")
	}
	oldDecision.Result = newState
	oldDecision.EmployeeID = &empID
	ds.makeDecision(oldDecision, reason)
	return ds.decisionRepository.Update(oldDecision)
}
*/
