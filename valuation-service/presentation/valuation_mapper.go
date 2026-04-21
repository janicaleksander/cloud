package presentation

import (
	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/valuationservice/application/command"
	"github.com/janicaleksander/cloud/valuationservice/application/query"
)

func GetValuationHTTPToQuery(id uuid.UUID) *query.GetValuationQuery {
	return &query.GetValuationQuery{
		ValuationID: id.String(),
	}
}

func GetValuationsHTTPToQuery() *query.GetValuationsQuery {
	return &query.GetValuationsQuery{}
}

func DeleteValuationHTTPToCommand(id uuid.UUID) *command.DeleteValuationCommand {
	return &command.DeleteValuationCommand{
		ID: id.String(),
	}

}
