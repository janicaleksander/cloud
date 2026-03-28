package application

import (
	"context"

	"github.com/janicaleksander/cloud/claimservice/domain"
	"github.com/janicaleksander/cloud/claimservice/persistance"
)

type ClaimService struct {
	claimRepository *persistance.ClaimRepository
	//probably to rabbit
}

func NewClaimService(claimRepo *persistance.ClaimRepository) *ClaimService {
	return &ClaimService{claimRepository: claimRepo}
}

//http methods

func (c *ClaimService) GetClaim(id uint) (*domain.Claim, error) {
	return c.claimRepository.GetById(context.Background(), id)
}
func (c *ClaimService) GetClaims() {}
func (c *ClaimService) CreateClaim(claim *domain.Claim) error {
	claim.Status = domain.NEW
	return c.claimRepository.Save(context.Background(), claim)
}
func (c *ClaimService) DeleteClaim() {}
func (c *ClaimService) UpdateClim()  {}

//rabbit events methods

func (c *ClaimService) ReceiveClaimEvent() {
	//receive
	//save to db
	//publish event
}

func (c *ClaimService) ChangeClaimStatus() {}
