package application

import (
	"context"
	"log/slog"

	"github.com/janicaleksander/cloud/common/event"
	"github.com/janicaleksander/cloud/valuationservice/domain"
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

func NewValuationService(valuationRepo domain.ValuationRepository, publisher ValuationPublisher, damageDetector DamageDetector) *ValuationService {
	slog.Info("Creating ValuationService")
	return &ValuationService{
		valuationRepository: valuationRepo,
		publisher:           publisher,
		damageDetector:      damageDetector,
	}
}
func (vs *ValuationService) CreateValuation(dV *domain.Valuation) (*domain.Valuation, error) {
	slog.Info("Creating valuation with ID", "claimID", dV.ClaimID)
	return vs.valuationRepository.Save(context.Background(), dV)
}

func (vs *ValuationService) GetValuations() ([]*domain.Valuation, error) {
	slog.Info("Getting all valuations")
	return vs.valuationRepository.GetAll(context.Background())
}

func (vs *ValuationService) GetValuation(claimID uint) (*domain.Valuation, error) {
	slog.Info("Getting valuation for claimID", "claimID", claimID)
	return vs.valuationRepository.GetById(context.Background(), claimID)
}

func (vs *ValuationService) UpdateValuation(oldValuation *domain.Valuation, amount float64) (*domain.Valuation, error) {
	slog.Info("Updating valuation for claimID", "claimID")
	updated := *oldValuation
	if updated.Amount != amount && amount != 0 {
		updated.Amount = amount
	}
	return vs.valuationRepository.Update(context.Background(), &updated)
}

func (vs *ValuationService) DeleteValuation(valuationID uint) error {
	slog.Info("Deleting valuation with ID", "claimID", valuationID)
	return vs.valuationRepository.DeleteById(context.Background(), valuationID)
}

func (vs *ValuationService) CalculateValuation(urls []string, claimID uint) error {
	slog.Info("Calculating valuation for claimID", "claimID", claimID)
	existing, err := vs.valuationRepository.GetById(context.Background(), claimID)
	if err == nil && existing != nil {
		return vs.publisher.Publish("events", event.ValuationCalculatedEvent{
			ClaimID:      existing.ClaimID,
			PayoutAmount: existing.Amount,
		})
	}
	damages, err := vs.damageDetector.Analyze(context.Background(), urls)
	if err != nil {
		return err
	}
	parts := make([]*domain.Part, len(damages))
	for i, damage := range damages {
		parts[i] = &domain.Part{
			Name: damage,
			Cost: float64(len(damages)) * 1000, //this is mock
		}
	}
	//this is mock
	amount := 0.0
	for _, part := range parts {
		amount += part.Cost
	}
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
	})
}
