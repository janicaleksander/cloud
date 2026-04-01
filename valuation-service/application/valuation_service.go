package application

import (
	"context"

	"github.com/janicaleksander/cloud/common/event"
	"github.com/janicaleksander/cloud/valuationservice/domain"
	"github.com/janicaleksander/cloud/valuationservice/persistance"
)

type ValuationService struct {
	valuationRepository domain.ValuationRepository
	publisher           ValuationPublisher
	damageDetector      DamageDetector
}

type ValuationPublisher interface {
	Publish(exchange string, msg interface{}) error
}

type DamageDetector interface {
	Analyze(ctx context.Context, urls []string) ([]string, error)
}

func NewValuationService(valuationRepo *persistance.ValuationRepository, publisher ValuationPublisher, damageDetector DamageDetector) *ValuationService {
	return &ValuationService{
		valuationRepository: valuationRepo,
		publisher:           publisher,
		damageDetector:      damageDetector,
	}
}
func (vs *ValuationService) CreateValuation(dV *domain.Valuation) (*domain.Valuation, error) {
	return vs.valuationRepository.Save(context.Background(), dV)
}

func (vs *ValuationService) GetValuations() ([]*domain.Valuation, error) {
	return vs.valuationRepository.GetAll(context.Background())
}

func (vs *ValuationService) GetValuation(claimID uint) (*domain.Valuation, error) {
	return vs.valuationRepository.GetById(context.Background(), claimID)
}

func (vs *ValuationService) UpdateValuation(oldValuation *domain.Valuation, amount float64) (*domain.Valuation, error) {
	updated := *oldValuation
	if updated.Amount != amount && amount != 0 {
		updated.Amount = amount
	}
	return vs.valuationRepository.Update(context.Background(), &updated)
}

func (vs *ValuationService) DeleteValuation(valuationID uint) error {
	return vs.valuationRepository.DeleteById(context.Background(), valuationID)
}

func (vs *ValuationService) CalculateValuation(urls []string, claimID uint) error {

	existing, err := vs.valuationRepository.GetById(context.Background(), claimID)
	if err == nil && existing != nil {
		return nil
	}

	damages, err := vs.damageDetector.Analyze(context.Background(), urls)
	if err != nil {
		return err
	}
	parts := make([]*domain.Part, len(damages))
	for i, damage := range damages {
		parts[i] = &domain.Part{
			Name: damage,
			Cost: 1000, //this is mock
		}
	}
	//this is mock
	amount := float64(len(damages)) * 1000
	_, err = vs.CreateValuation(&domain.Valuation{
		ClaimID: claimID,
		Amount:  amount,
		Parts:   parts,
	})
	if err != nil {
		return err
	}

	return vs.publisher.Publish("events", event.ValuationCalculatedEvent{
		ClaimID:      claimID,
		PayoutAmount: amount,
		DamageItems:  damages,
	})
}

/*
TODO 📌 Inny problem architektoniczny
func NewValuationService(valuationRepo *persistance.ValuationRepository, ...)


❌ zależysz od konkretnej implementacji

Powinno być:

func NewValuationService(valuationRepo domain.ValuationRepository, ...)


👉 bo:

application → zależy od interfejsów
infrastructure → implementuje

*/
