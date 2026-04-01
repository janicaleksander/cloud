package application

import (
	"context"

	"github.com/janicaleksander/cloud/claimservice/domain"
	"github.com/janicaleksander/cloud/claimservice/persistance"
	"github.com/janicaleksander/cloud/common/event"
)

type ClaimService struct {
	claimRepository domain.ClaimRepository
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
	claim.Status = domain.NEW

	urls := make([]string, 0, len(claim.Files))
	for idx := range claim.Files {
		urls = append(urls, claim.Files[idx].StorageURL)
	}
	savedClaim, err := c.claimRepository.Save(context.Background(), claim)
	if err != nil {
		return err
	}
	err = c.publisher.Publish("events", event.RegisterUserForNotificationEvent{
		ClaimID: savedClaim.ID,
		Email:   savedClaim.Email,
	})
	if err != nil {
		return err
	}

	err = c.pushClaimSubmittedEvent(&event.ClaimSubmittedEvent{
		UserID:       savedClaim.UserID,
		ClaimID:      savedClaim.ID,
		VIN:          savedClaim.VIN,
		AccidentDate: savedClaim.AccidentDate,
		StorageURL:   urls,
	})

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

func (c *ClaimService) UpdateClaim(oldClaimDomain *domain.Claim, newUserEmail string) error {
	if newUserEmail != oldClaimDomain.Email && newUserEmail != "" {
		oldClaimDomain.Email = newUserEmail
	} //todo publish msg to change email in notification service

	_, err := c.claimRepository.Update(context.Background(), oldClaimDomain)
	if err != nil {
		return err
	}
	err = c.publisher.Publish("events", event.ChangeEmailForNotification{
		ClaimID: oldClaimDomain.ID,
		Email:   newUserEmail,
	})
	if err != nil {
		return err
	}
	return nil
}

//rabbit events methods

func (c *ClaimService) ChangeClaimStatus(claimID uint, newStatus domain.Status) error {
	return c.claimRepository.UpdateStatus(context.Background(), claimID, newStatus)
}
