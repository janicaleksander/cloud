package application

import (
	"context"
	"log/slog"

	"github.com/janicaleksander/cloud/claimservice/domain"
	"github.com/janicaleksander/cloud/claimservice/persistence"
	"github.com/janicaleksander/cloud/common/event"
)

type ClaimService struct {
	claimRepository domain.ClaimRepository // //under the hood there is ref do persistance
	publisher       ClaimEventPublisher
}

type ClaimEventPublisher interface {
	Publish(exchange string, msg interface{}) error
}

func NewClaimService(claimRepo *persistence.ClaimRepository, publisher ClaimEventPublisher) *ClaimService {
	slog.Info("Creating ClaimService")
	return &ClaimService{
		claimRepository: claimRepo,
		publisher:       publisher,
	}
}

//http methods

func (c *ClaimService) CreateClaim(claim *domain.Claim, domainFiles []*domain.File) (*domain.Claim, error) {
	slog.Info("Creating claim with ID: ", "claimID", claim.ID)
	claim.Status = domain.NEW
	for idx := range domainFiles {
		domainFiles[idx].StorageURL = "https://storage.example.com/" + domainFiles[idx].FileName
	}
	if len(domainFiles) != 0 {
		claim.Files = domainFiles
	}

	urls := make([]string, 0, len(claim.Files))
	for idx := range claim.Files {
		urls = append(urls, claim.Files[idx].StorageURL)
	}
	savedClaim, err := c.claimRepository.Save(context.Background(), claim)
	if err != nil {
		return nil, err
	}
	err = c.publisher.Publish("events", event.RegisterUserForNotificationEvent{
		ClaimID: savedClaim.ID,
		Email:   savedClaim.Email,
	})
	if err != nil {
		return nil, err
	}

	err = c.pushClaimSubmittedEvent(&event.ClaimSubmittedEvent{
		UserID:       savedClaim.UserID,
		ClaimID:      savedClaim.ID,
		VIN:          savedClaim.VIN,
		AccidentDate: savedClaim.AccidentDate,
		StorageURL:   urls,
	})

	if err != nil {
		return nil, err
	}
	return savedClaim, nil
}
func (c *ClaimService) pushClaimSubmittedEvent(e *event.ClaimSubmittedEvent) error {
	slog.Info("Publishing ClaimSubmittedEvent for ClaimID: ", "claimID", e.ClaimID)
	return c.publisher.Publish("events", *e)
}

func (c *ClaimService) GetClaim(id uint) (*domain.Claim, error) {
	slog.Info("Getting claim with ID: ", "claimID", id)
	return c.claimRepository.GetById(context.Background(), id)
}
func (c *ClaimService) GetClaims() ([]*domain.Claim, error) {
	slog.Info("Getting all claims")
	return c.claimRepository.GetAll(context.Background())
}

func (c *ClaimService) DeleteClaim(id uint) error {
	slog.Info("Deleting claim with ID: ", "claimID", id)
	return c.claimRepository.DeleteById(context.Background(), id)
}

func (c *ClaimService) UpdateClaim(oldClaimDomain *domain.Claim, newUserEmail string) (*domain.Claim, error) {
	slog.Info("Updating claim with ID: ", "claimID", oldClaimDomain.ID)
	if newUserEmail != oldClaimDomain.Email && newUserEmail != "" {
		oldClaimDomain.Email = newUserEmail
	}

	updatedClaim, err := c.claimRepository.Update(context.Background(), oldClaimDomain)
	if err != nil {
		return nil, err
	}
	err = c.publisher.Publish("events", event.ChangeEmailForNotification{
		ClaimID: oldClaimDomain.ID,
		Email:   newUserEmail,
	})
	if err != nil {
		return nil, err
	}
	return updatedClaim, nil
}

//rabbit events methods

func (c *ClaimService) ChangeClaimStatus(claimID uint, newStatus domain.Status) error {
	slog.Info("Changing claim status for ClaimID: ", "claimID", claimID, "newStatus", newStatus)
	return c.claimRepository.UpdateStatus(context.Background(), claimID, newStatus)
}
