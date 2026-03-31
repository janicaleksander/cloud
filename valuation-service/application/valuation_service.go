package application

import (
	"context"
	"math/rand"

	"github.com/janicaleksander/cloud/common/event"
	"github.com/janicaleksander/cloud/valuationservice/domain"
	"github.com/janicaleksander/cloud/valuationservice/persistance"
)

type ValuationService struct {
	valuationRepository domain.ValuationRepository
	publisher           ValuationPublisher
}

type ValuationPublisher interface {
	Publish(exchange string, msg interface{}) error
}

func NewValuationService(valuationRepo *persistance.ValuationRepository, publisher ValuationPublisher) *ValuationService {
	return &ValuationService{
		valuationRepository: valuationRepo,
		publisher:           publisher,
	}
}
func (vs *ValuationService) CreateValuation(claimID uint, amount float64) (*domain.Valuation, error) {
	valuation := &domain.Valuation{
		ClaimID: claimID,
		Amount:  amount,
	}
	return vs.valuationRepository.Save(context.Background(), valuation)
}

func (vs *ValuationService) GetValuations() ([]*domain.Valuation, error) {
	return vs.valuationRepository.GetAll(context.Background())
}

func (vs *ValuationService) GetValuation(claimID uint) (*domain.Valuation, error) {
	return vs.valuationRepository.GetById(context.Background(), claimID)
}

func (vs *ValuationService) UpdateValuation(oldValuation *domain.Valuation, amount float64) (*domain.Valuation, error) {
	if oldValuation.Amount != amount && amount != 0 {
		oldValuation.Amount = amount
	}
	return vs.valuationRepository.Update(context.Background(), oldValuation)
}

func (vs *ValuationService) DeleteValuation(valuationID uint) error {
	return vs.valuationRepository.DeleteById(context.Background(), valuationID)
}

func (vs *ValuationService) CalculateValuation(urls []string, claimID uint) {

	// Mock valuation: random amount between 500 and 10000
	amount := rand.Float64()*(10000-500) + 500

	// Generate random parts for the valuation
	parts := []string{"bumper", "hood", "door", "fender", "headlight"}
	randomParts := make([]string, rand.Intn(3)+1) // 1 to 3 random parts
	for i := range randomParts {
		randomParts[i] = parts[rand.Intn(len(parts))]
	}
	_, err := vs.CreateValuation(claimID, amount)
	if err != nil {
		//TODO logs
	}
	err = vs.publisher.Publish("events", event.ValuationCalculatedEvent{
		ClaimID:      claimID,
		PayoutAmount: amount,
		DamageItems:  parts,
	})
	if err != nil {
		//toodlogs
	}

}
