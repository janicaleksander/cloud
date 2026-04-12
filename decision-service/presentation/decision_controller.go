package presentation

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/janicaleksander/cloud/decisionservice/application"
)

type DecisionController struct {
	decisionService *application.DecisionService
}

func NewDecisionController(decisionService *application.DecisionService) *DecisionController {
	slog.Info("Creating DecisionController")
	return &DecisionController{
		decisionService: decisionService,
	}
}

func success(w http.ResponseWriter, msg any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(msg)
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
	domainDecisions, err := d.decisionService.GetDecisions()
	if err != nil {
		failure(w, http.StatusInternalServerError, err.Error())
		return
	}
	dto := make([]*GetDecisionResponseDTO, 0, len(domainDecisions))
	for _, decision := range domainDecisions {
		dto = append(dto, GetDecisionDomainToResponse(decision))
	}
	success(w, map[string]any{"decisions": dto})
}

func (d *DecisionController) GetDecisionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		failure(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	slog.Info("HTTP GetDecisionHandler called")
	decisionID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		failure(w, http.StatusBadRequest, "Invalid decision ID")
		return
	}
	domainDecision, err := d.decisionService.GetDecision(uint(decisionID))
	if err != nil {
		failure(w, http.StatusNotFound, err.Error())
		return
	}
	dto := GetDecisionDomainToResponse(domainDecision)
	success(w, map[string]any{"decision": dto})
}

func (d *DecisionController) GetWaitingDecisionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		failure(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	slog.Info("HTTP GetWaitingDecisionsHandler called")
	domainDecisions, err := d.decisionService.GetWaitingDecisions()
	if err != nil {
		failure(w, http.StatusInternalServerError, err.Error())
		return
	}
	dto := make([]*GetDecisionResponseDTO, 0, len(domainDecisions))
	for _, decision := range domainDecisions {
		dto = append(dto, GetDecisionDomainToResponse(decision))
	}
	success(w, map[string]any{"decisions": dto})

}

func (d *DecisionController) UpdateDecisionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		failure(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	slog.Info("HTTP UpdateDecisionHandler called")
	decisionID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		failure(w, http.StatusBadRequest, "Invalid decision ID")
		return
	}
	var dto UpdateDecisionRequestDTO
	err = json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		failure(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	oldDecisionDomain, err := d.decisionService.GetDecision(uint(decisionID))
	if err != nil {
		failure(w, http.StatusNotFound, err.Error())
		return
	}
	updatedDecision, err := d.decisionService.UpdateDecisionState(oldDecisionDomain, dto.NewState, dto.EmpID, dto.Reason)
	if err != nil {
		failure(w, http.StatusInternalServerError, err.Error())
		return
	}
	success(w, map[string]any{"decision": GetDecisionDomainToResponse(updatedDecision)})

}

func (d *DecisionController) DeleteDecisionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		failure(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	slog.Info("HTTP DeleteDecisionHandler called")
	decisionID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		failure(w, http.StatusBadRequest, "Invalid decision ID")
		return
	}
	err = d.decisionService.DeleteDecision(uint(decisionID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}
	success(w, map[string]string{"message": "Decision deleted successfully + id: " + strconv.Itoa(decisionID)})
}
