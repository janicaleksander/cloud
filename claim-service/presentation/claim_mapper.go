package presentation

import (
	"os"

	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/claimservice/application/command"
	"github.com/janicaleksander/cloud/claimservice/application/query"
)

func CreateClaimHTTPToCommand(newID uuid.UUID, dto *CreateClaimRequestDTO, files []*os.File) *command.CreateClaimCommand {
	return &command.CreateClaimCommand{
		ID:           newID.String(),
		UserID:       dto.UserID,
		Email:        dto.Email,
		VIN:          dto.VIN,
		AccidentDate: dto.AccidentDate,
		ObjectFiles:  files,
	}
}

func GetClaimHTTPToQuery(claimID string) *query.GetClaimByIdQuery {
	return &query.GetClaimByIdQuery{ClaimID: claimID}
}

func GetClaimsHTTPToQuery() *query.GetClaimsQuery {
	return &query.GetClaimsQuery{}
}

func DeleteClaimHTTPToCommand(claimID string) *command.DeleteClaimCommand {
	return &command.DeleteClaimCommand{ClaimID: claimID}
}

func GetFileHTTPToQuery(fileID string) *query.GetFileFromStorageQuery {
	return &query.GetFileFromStorageQuery{FileID: fileID}
}
