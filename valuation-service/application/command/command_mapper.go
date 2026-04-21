package command

import (
	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/valuationservice/domain"
)

func CreatePartDomainToCommand(p *domain.Part) *PartCommand {
	return &PartCommand{
		ID:   p.ID.String(),
		Name: p.Name,
		Cost: p.Cost,
	}
}

func CreateValuationDomainToCommand(d *domain.Valuation) *CreateValuationCommand {
	parts := make([]PartCommand, len(d.Parts))
	for i, part := range d.Parts {
		parts[i] = *CreatePartDomainToCommand(part)
	}
	return &CreateValuationCommand{
		ID:      d.ID.String(),
		ClaimID: d.ClaimID.String(),
		Amount:  d.Amount,
		Parts:   parts,
	}
}

func CretePartCommandToDomain(p *PartCommand) *domain.Part {
	id, err := uuid.Parse(p.ID)
	if err != nil {
		id = uuid.New()
	}
	return &domain.Part{
		ID:   id,
		Name: p.Name,
		Cost: p.Cost,
	}
}

func CreateValuationCommandToDomain(c *CreateValuationCommand) *domain.Valuation {
	id, err := uuid.Parse(c.ID)
	if err != nil {
		id = uuid.New()
	}
	claimID, err := uuid.Parse(c.ClaimID)
	if err != nil {
		claimID = uuid.New()
	}
	parts := make([]*domain.Part, len(c.Parts))
	for i, part := range c.Parts {
		parts[i] = CretePartCommandToDomain(&part)
	}
	return &domain.Valuation{
		ID:      id,
		ClaimID: claimID,
		Amount:  c.Amount,
		Parts:   parts,
	}

}
