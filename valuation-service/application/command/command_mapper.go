package command

import (
	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/valuationservice/domain"
)

func CreateValuationCommandToDomain(cmd *CreateValuationCommand) *domain.Valuation {
	parts := make([]*domain.Part, 0, len(cmd.Parts))

	for idx := range cmd.Parts {
		pid, err := uuid.Parse(cmd.Parts[idx].ID)
		if err != nil {
			return nil
		}
		parts = append(parts, &domain.Part{
			ID:   pid,
			Name: cmd.Parts[idx].Name,
			Cost: cmd.Parts[idx].Cost,
		})
	}
	vid, err := uuid.Parse(cmd.ID)
	if err != nil {
		return nil
	}
	cid, err := uuid.Parse(cmd.ClaimID)
	if err != nil {
		return nil
	}
	return &domain.Valuation{
		ID:      vid,
		ClaimID: cid,
		Amount:  cmd.Amount,
		Parts:   parts,
	}

}

func CreateValuationDomainToCommand(val *domain.Valuation) *CreateValuationCommand {
	parts := make([]PartCommand, 0, len(val.Parts))

	for idx := range val.Parts {
		parts = append(parts, PartCommand{
			ID:   val.Parts[idx].ID.String(),
			Name: val.Parts[idx].Name,
			Cost: val.Parts[idx].Cost,
		})
	}
	return &CreateValuationCommand{
		ID:      val.ID.String(),
		ClaimID: val.ClaimID.String(),
		Amount:  val.Amount,
		Parts:   parts,
	}

}
