package presentation

import (
	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/decisionservice/application/command"
	"github.com/janicaleksander/cloud/decisionservice/application/query"
)

func GetDecisionHTTPToQuery(decisionID uuid.UUID) *query.GetDecisionQuery {
	return &query.GetDecisionQuery{DecisionID: decisionID.String()}
}

func GetDecisionsHTTPToQuery() *query.GetDecisionsQuery {
	return &query.GetDecisionsQuery{}
}

func GetWaitingDecisionsHTTPToQuery() *query.GetWaitingDecisionsQuery {
	return &query.GetWaitingDecisionsQuery{}
}

func DeleteDecisionHTTPToCommand(decisionID uuid.UUID) *command.DeleteDecisionCommand {
	return &command.DeleteDecisionCommand{DecisionID: decisionID.String()}
}

func UpdateDecisionHTTPToCommand(decisionID uuid.UUID, req *UpdateDecisionRequestDTO) *command.UpdateDecisionStateCommand {
	return &command.UpdateDecisionStateCommand{
		DecisionID: decisionID.String(),
		NewState:   req.NewState,
		EmpID:      req.EmpID,
		Reason:     req.Reason,
	}
}
