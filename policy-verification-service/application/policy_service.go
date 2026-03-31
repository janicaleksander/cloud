package application

import (
	"context"
	"log"
	"time"

	"github.com/janicaleksander/cloud/common/event"
	"github.com/janicaleksander/cloud/policyverificationservice/domain"
	"github.com/janicaleksander/cloud/policyverificationservice/persistance"
)

type PolicyService struct {
	policyRepository domain.PolicyRepository
	publisher        PolicyEventPublisher
}

type PolicyEventPublisher interface {
	Publish(exchange string, msg interface{}) error
}

func NewPolicyService(policyRepository *persistance.PolicyRepository, publisher PolicyEventPublisher) *PolicyService {
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
	updated := *newPolicy

	if newFrom != (time.Time{}) {
		updated.From = newFrom
	}
	if newTo != (time.Time{}) {
		updated.To = newTo
	}
	_, err := s.policyRepository.Update(context.Background(), &updated)
	return err
}
func (s *PolicyService) DeletePolicy(policyID uint) error {
	return s.policyRepository.DeleteById(context.Background(), policyID)
}

func (s *PolicyService) CheckUserPolicy(claimID uint, userID uint, vin string, accidentDate time.Time, urls []string) {
	hasPolicy, policy := s.policyRepository.IfUserHasPolicy(context.Background(), userID, vin)

	if !hasPolicy {
		err := s.publisher.Publish("events", event.PolicyDeniedEvent{
			ClaimID: claimID,
			Reason:  string(domain.PolicyNotFound),
		})
		if err != nil {
			log.Printf("Failed to publish PolicyDeniedEvent for claimID %d: %v", claimID, err)
		}
		return
	}

	valid, reason := policy.IsValid(accidentDate)

	if valid {
		err := s.publisher.Publish("events", event.PolicyVerifiedEvent{
			ClaimID:    claimID,
			StorageURL: urls,
		})
		if err != nil {
			log.Printf("Failed to publish PolicyVerifiedEvent for claimID %d: %v", claimID, err)
		}
	} else {
		err := s.publisher.Publish("events", event.PolicyDeniedEvent{
			ClaimID: claimID,
			Reason:  string(reason),
		})
		if err != nil {
			log.Printf("Failed to publish PolicyDeniedEvent for claimID %d: %v", claimID, err)
		}
	}
}
