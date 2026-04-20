package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/common/event"
	"github.com/janicaleksander/cloud/valuationservice/application/interfaces"
	"github.com/janicaleksander/cloud/valuationservice/domain"
	"github.com/mehdihadeli/go-mediatr"
)

type CalculateValuationCommand struct {
	ClaimID string
	Urls    []string
}

type CalculateValuationCommandHandler struct {
	repo      domain.ValuationRepository
	detector  interfaces.DamageDetector
	publisher interfaces.ValuationPublisher
}

func NewCalculateValuationCommandHandler(repo domain.ValuationRepository, detector interfaces.DamageDetector, publisher interfaces.ValuationPublisher) *CalculateValuationCommandHandler {
	return &CalculateValuationCommandHandler{
		repo:      repo,
		detector:  detector,
		publisher: publisher,
	}
}
func (h *CalculateValuationCommandHandler) SelfRegister() error {
	return mediatr.RegisterRequestHandler[*CalculateValuationCommand, *mediatr.Unit](h)
}

func (h *CalculateValuationCommandHandler) Handle(ctx context.Context, cmd *CalculateValuationCommand) (*mediatr.Unit, error) {
	cid, err := uuid.Parse(cmd.ClaimID)
	if err != nil {
		return nil, err
	}
	existing, err := h.repo.GetById(context.Background(), cid)
	if err == nil && existing != nil {
		err := h.publisher.Publish("events", event.ValuationCalculatedEvent{
			ClaimID:      cmd.ClaimID,
			PayoutAmount: existing.Amount,
		})
		if err != nil {
			return nil, err
		}
	}

	damages, err := h.detector.Analyze(context.Background(), cmd.Urls)
	if err != nil {
		return nil, err
	}
	parts := make([]*domain.Part, len(damages))
	for i, damage := range damages {
		parts[i] = &domain.Part{
			ID:   uuid.New(),
			Name: damage,
			Cost: float64(len(damages)) * 1000, //this is mock
		}
	}
	//this is mock
	amount := 0.0
	for _, part := range parts {
		amount += part.Cost
	}

	//send
	valuationDomain := &domain.Valuation{
		ID:      uuid.New(),
		ClaimID: cid,
		Amount:  amount,
		Parts:   parts,
	}
	createValuationCmd := CreateValuationDomainToCommand(valuationDomain)
	_, err = mediatr.Send[*CreateValuationCommand, *mediatr.Unit](ctx, createValuationCmd)

	err = h.publisher.Publish("events", event.ValuationCalculatedEvent{
		ClaimID:      cmd.ClaimID,
		PayoutAmount: amount,
	})
	if err != nil {
		return nil, err
	}
	return &mediatr.Unit{}, nil
}
