package application

import (
	"context"
	"log/slog"
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
	slog.Info("Creating PolicyService with provided PolicyRepository and PolicyEventPublisher")
	return &PolicyService{
		policyRepository: policyRepository,
		publisher:        publisher,
	}
}

func (s *PolicyService) CreatePolicy(policy *domain.Policy) (*domain.Policy, error) {
	slog.Info("Created policy with UserID: ", "userID", policy.UserID, "vin", policy.VIN)
	savedPolicy, err := s.policyRepository.Save(context.Background(), policy)
	return savedPolicy, err
}

func (s *PolicyService) GetPolicy(policyId uint) (*domain.Policy, error) {
	slog.Info("Getting policy by ID", "policyId", policyId)
	return s.policyRepository.GetById(context.Background(), policyId)
}

func (s *PolicyService) GetPolicies() ([]*domain.Policy, error) {
	slog.Info("Getting all policies")
	return s.policyRepository.GetAll(context.Background())
}

func (s *PolicyService) UpdatePolicy(newPolicy *domain.Policy, newFrom, newTo time.Time) (*domain.Policy, error) {
	slog.Info("Updating policy with ID: ", "policyId", newPolicy.ID)
	updated := *newPolicy

	if newFrom != (time.Time{}) {
		updated.From = newFrom
	}
	if newTo != (time.Time{}) {
		updated.To = newTo
	}
	updatedPolicy, err := s.policyRepository.Update(context.Background(), &updated)
	return updatedPolicy, err
}
func (s *PolicyService) DeletePolicy(policyID uint) error {
	slog.Info("Deleting policy with ID: ", "policyId", policyID)
	return s.policyRepository.DeleteById(context.Background(), policyID)
}

func (s *PolicyService) CheckUserPolicy(claimID uint, userID uint, vin string, accidentDate time.Time, urls []string) {
	slog.Info("Checking user policy", "claimID", claimID, "userID", userID, "vin", vin)
	hasPolicy, policy := s.policyRepository.IfUserHasPolicy(context.Background(), userID, vin)

	if !hasPolicy {
		err := s.publisher.Publish("events", event.PolicyDeniedEvent{
			ClaimID: claimID,
			Reason:  string(domain.PolicyNotFound),
		})
		if err != nil {
			slog.Error("Failed to publish PolicyDeniedEvent for claimID", "claimID", claimID, "error", err)
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
			slog.Error("Failed to publish PolicyVerifiedEvent for claimID", "claimID", claimID, "error", err)
		}
	} else {
		err := s.publisher.Publish("events", event.PolicyDeniedEvent{
			ClaimID: claimID,
			Reason:  string(reason),
		})
		if err != nil {
			slog.Error("Failed to publish PolicyDeniedEvent for claimID", "claimID", claimID, "error", err)
		}
	}
}
