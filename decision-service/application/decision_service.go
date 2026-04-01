package application

import (
	"errors"
	"fmt"
	"log"

	"github.com/janicaleksander/cloud/common/event"
	"github.com/janicaleksander/cloud/decisionservice/domain"
)

type DecisionService struct {
	decisionRepository domain.DecisionRepository
	publisher          DecisionPublisher
}

type DecisionPublisher interface {
	Publish(exchange string, message any) error
}

func NewDecisionService(decisionRepo domain.DecisionRepository, publisher DecisionPublisher) *DecisionService {
	return &DecisionService{
		decisionRepository: decisionRepo,
		publisher:          publisher,
	}
}

func (ds *DecisionService) makeDecision(newDecision *domain.Decision, reason string) {
	fmt.Println(newDecision.Result)
	if newDecision.Result == domain.ACCEPTED {
		err := ds.publisher.Publish("events", event.PayoutApprovedEvent{
			ClaimID:              newDecision.ClaimID,
			AcceptedPayoutAmount: newDecision.Payout,
		})
		fmt.Println("wyslalem na publish")

		if err != nil {
			log.Println("Failed to publish PayoutApprovedEvent:", err)
		}
	}
	if newDecision.Result == domain.REJECTED {
		err := ds.publisher.Publish("events", event.PayoutRejectedEvent{
			ClaimID: newDecision.ClaimID,
			Reason:  reason,
		})
		fmt.Println("wyslalem na rjeected")

		if err != nil {
			log.Println("Failed to publish PayoutRejectedEvent:", err)
		}
	}
}
func (ds *DecisionService) PrepareDecision(claimID uint, payoutAmount float64) (*domain.Decision, error) {
	decisionDomain := &domain.Decision{
		ClaimID: claimID,
		Payout:  payoutAmount,
		Result:  domain.WAITING,
	}
	fmt.Println("jestem")
	createdDecision, err := ds.decisionRepository.Save(decisionDomain)
	if err != nil {
		return nil, err
	}
	return createdDecision, nil
}

func (ds *DecisionService) GetDecision(decisionID uint) (*domain.Decision, error) {
	return ds.decisionRepository.GetByID(decisionID)
}
func (ds *DecisionService) GetDecisions() ([]*domain.Decision, error) {
	return ds.decisionRepository.GetAll()

}
func (ds *DecisionService) GetWaitingDecisions() ([]*domain.Decision, error) {
	return ds.decisionRepository.GetAllWaiting()
}
func (ds *DecisionService) DeleteDecision(decisionID uint) error {
	return ds.decisionRepository.DeleteById(decisionID)

}
func (ds *DecisionService) UpdateDecisionState(oldDecision *domain.Decision, newState domain.DecisionResult, empID uint, reason string) (*domain.Decision, error) {
	if oldDecision.Result != domain.WAITING {
		return nil, errors.New("already accepted/denied")
	}
	//send rabbit mq accept/deny
	oldDecision.Result = newState
	oldDecision.EmployeeID = &empID
	ds.makeDecision(oldDecision, reason)
	return ds.decisionRepository.Update(oldDecision)
}
