package presentation

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/janicaleksander/cloud/decisionservice/application/command"
	"github.com/janicaleksander/cloud/decisionservice/application/query"
	"github.com/mehdihadeli/go-mediatr"
)

type DecisionController struct{}

func NewDecisionController() *DecisionController {
	slog.Info("Creating DecisionController")
	return &DecisionController{}
}

func success(w http.ResponseWriter, msg any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if msg != nil {
		json.NewEncoder(w).Encode(msg)
	}
}

func successWithLocation(w http.ResponseWriter, location string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", location)
	w.WriteHeader(code)
}

func failure(w http.ResponseWriter, statusCode int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func (d *DecisionController) GetDecisionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		failure(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	slog.Info("HTTP GetDecisionsHandler called")
	q := &query.GetDecisionsQuery{}
	decisions, err := mediatr.Send[*query.GetDecisionsQuery, *query.GetDecisionsQueryResult](context.Background(), q)
	if err != nil {
		failure(w, http.StatusInternalServerError, err.Error())
		return
	}
	success(w, map[string]any{"decisions": decisions.Decisions}, 200)
}

func (d *DecisionController) GetDecisionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		failure(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	slog.Info("HTTP GetDecisionHandler called")
	decisionID := chi.URLParam(r, "id")

	q := &query.GetDecisionQuery{DecisionID: decisionID}

	decision, err := mediatr.Send[*query.GetDecisionQuery, *query.GetDecisionQueryResult](context.Background(), q)
	if err != nil {
		failure(w, http.StatusNotFound, err.Error())
		return
	}
	success(w, map[string]any{"decision": decision}, 200)
}

func (d *DecisionController) GetWaitingDecisionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		failure(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	slog.Info("HTTP GetWaitingDecisionsHandler called")
	q := &query.GetWaitingDecisionsQuery{}
	waitingDecisions, err := mediatr.Send[*query.GetWaitingDecisionsQuery, *query.GetWaitingDecisionsQueryResponse](context.Background(), q)
	if err != nil {
		failure(w, http.StatusInternalServerError, err.Error())
		return
	}
	success(w, map[string]any{"decisions": waitingDecisions.Waiting}, 200)

}

func (d *DecisionController) DeleteDecisionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		failure(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	slog.Info("HTTP DeleteDecisionHandler called")
	decisionID := chi.URLParam(r, "id")

	cmd := &command.DeleteDecisionCommand{DecisionID: decisionID}

	_, err := mediatr.Send[*command.DeleteDecisionCommand, *mediatr.Unit](context.Background(), cmd)

	if err != nil {
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}
	success(w, nil, http.StatusNoContent)
}

func (d *DecisionController) UpdateDecisionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		failure(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	slog.Info("HTTP UpdateDecisionHandler called")
	decisionID := chi.URLParam(r, "id")

	var dto UpdateDecisionRequestDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		failure(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	cmd := &command.UpdateDecisionStateCommand{
		DecisionID: decisionID,
		NewState:   dto.NewState,
		EmpID:      dto.EmpID,
	}
	_, err = mediatr.Send[*command.UpdateDecisionStateCommand, *mediatr.Unit](context.Background(), cmd)
	if err != nil {
		failure(w, http.StatusNotFound, "Decision not found")
		return
	}
	successWithLocation(w, "/decisions/"+decisionID, http.StatusOK)
}
