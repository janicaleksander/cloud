package application

import (
	"context"

	"github.com/janicaleksander/cloud/claimservice/domain"
	"github.com/janicaleksander/cloud/claimservice/persistance"
	"github.com/janicaleksander/cloud/common/event"
)

type ClaimService struct {
	claimRepository domain.ClaimerRepository
	publisher       ClaimEventPublisher //under the hood there is ref do persistance
}

type ClaimEventPublisher interface {
	Publish(exchange string, msg interface{}) error
}

func NewClaimService(claimRepo *persistance.ClaimRepository, publisher ClaimEventPublisher) *ClaimService {
	return &ClaimService{
		claimRepository: claimRepo,
		publisher:       publisher,
	}
}

//http methods

func (c *ClaimService) CreateClaim(claim *domain.Claim) error {
	urls := make([]string, 0, len(claim.Files))
	for idx := range claim.Files {
		urls = append(urls, claim.Files[idx].StorageURL)
	}
	err := c.pushClaimSubmittedEvent(&event.ClaimSubmittedEvent{
		UserID:     claim.UserID,
		ClaimID:    claim.ID,
		StorageURL: urls,
	})
	if err != nil {
		return err
	}
	err = c.claimRepository.Save(context.Background(), claim)
	if err != nil {
		return err
	}
	return nil
}
func (c *ClaimService) pushClaimSubmittedEvent(e *event.ClaimSubmittedEvent) error {
	return c.publisher.Publish("events", *e)
}

func (c *ClaimService) GetClaim(id uint) (*domain.Claim, error) {
	return c.claimRepository.GetById(context.Background(), id)
}
func (c *ClaimService) GetClaims() ([]*domain.Claim, error) {
	return c.claimRepository.GetAll(context.Background())
}

func (c *ClaimService) DeleteClaim(id uint) error {
	return c.claimRepository.DeleteById(context.Background(), id)
}
func (c *ClaimService) UpdateClaim(d *domain.Claim) error {
	return c.claimRepository.Update(context.Background(), d)
}

//rabbit events methods

func (c *ClaimService) ChangeClaimStatus(claimID uint, newStatus domain.Status) error {
	claim, err := c.GetClaim(claimID)
	if err != nil {
		return err
	}
	claim.Status = newStatus
	return c.UpdateClaim(claim)
}
