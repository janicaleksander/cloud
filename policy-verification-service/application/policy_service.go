package application

import (
	"context"
	"log"
	"time"

	"github.com/janicaleksander/cloud/common/event"
	"github.com/janicaleksander/cloud/policyverificationservice/domain"
)

type PolicyService struct {
	policyRepository domain.PolicyRepository
	publisher        PolicyEventPublisher
}

type PolicyEventPublisher interface {
	Publish(exchange string, msg interface{}) error
}

func NewPolicyService(policyRepository domain.PolicyRepository, publisher PolicyEventPublisher) *PolicyService {
	return &PolicyService{
		policyRepository: policyRepository,
		publisher:        publisher,
	}
}

func (s *PolicyService) CreatePolicy(policy *domain.Policy) error {
	_, err := s.policyRepository.Save(context.Background(), policy)
	return err
}

func (s *PolicyService) GetPolicy(policyId uint) (*domain.Policy, error) {
	return s.policyRepository.GetById(context.Background(), policyId)
}

func (s *PolicyService) GetPolicies() ([]*domain.Policy, error) {
	return s.policyRepository.GetAll(context.Background())
}

func (s *PolicyService) UpdatePolicy(newPolicy *domain.Policy, newFrom, newTo time.Time) error {
	if newFrom != (time.Time{}) {
		newPolicy.From = newFrom
	}
	if newTo != (time.Time{}) {
		newPolicy.To = newTo
	}
	_, err := s.policyRepository.Update(context.Background(), newPolicy)
	return err
}
func (s *PolicyService) DeletePolicy(policyID uint) error {
	return s.policyRepository.DeleteById(context.Background(), policyID)
}

func (s *PolicyService) CheckUserPolicy(claimID uint, userID uint, vin string, accidentDate time.Time, urls []string) {
	// if user has a policy
	hasPolicy, policy := s.policyRepository.IfUserHasPolicy(context.Background(), userID, vin)
	// if policy is not expired
	ok := hasPolicy && policy.IsValid(accidentDate)
	if ok {
		err := s.publisher.Publish("events", event.PolicyVerifiedEvent{
			ClaimID:    claimID,
			StorageURL: urls,
		}) //??? TODO

		if err != nil {
			log.Printf("Failed to publish PolicyVerifiedEvent for claimID %d: %v", claimID, err)
			//TODO log
		}

	} else {
		err := s.publisher.Publish("events", event.PolicyDeniedEvent{
			ClaimID: claimID,
			Reason:  "some reason", // TODO
		}) //??? TODO

		if err != nil {
			log.Println("Failed to publish PolicyDeniedEvent for claimID %d: %v", claimID, err)
			//TODO log
		}

	}
}
