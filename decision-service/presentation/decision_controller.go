package presentation

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/janicaleksander/cloud/decisionservice/application"
)

type DecisionController struct {
	decisionService *application.DecisionService
}

func NewDecisionController(decisionService *application.DecisionService) *DecisionController {
	return &DecisionController{
		decisionService: decisionService,
	}
}

func (d *DecisionController) GetDecisionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	domainDecisions, err := d.decisionService.GetDecisions()
	if err != nil {
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}
	dto := make([]*GetDecisionResponseDTO, 0, len(domainDecisions))
	for _, decision := range domainDecisions {
		dto = append(dto, GetDecisionDomainToResponse(decision))
	}
	json.NewEncoder(w).Encode(dto)
}

func (d *DecisionController) GetDecisionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	decisionID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}
	domainDecision, err := d.decisionService.GetDecision(uint(decisionID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}
	dto := GetDecisionDomainToResponse(domainDecision)
	json.NewEncoder(w).Encode(dto)
}

func (d *DecisionController) GetWaitingDecisionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	domainDecisions, err := d.decisionService.GetWaitingDecisions()
	if err != nil {
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}
	dto := make([]*GetDecisionResponseDTO, 0, len(domainDecisions))
	for _, decision := range domainDecisions {
		dto = append(dto, GetDecisionDomainToResponse(decision))
	}
	json.NewEncoder(w).Encode(dto)

}

func (d *DecisionController) UpdateDecisionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	decisionID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}
	var dto UpdateDecisionRequestDTO
	err = json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}
	oldDecisionDomain, err := d.decisionService.GetDecision(uint(decisionID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}
	_, err = d.decisionService.UpdateDecisionState(oldDecisionDomain, dto.NewState, dto.EmpID, dto.Reason)
	if err != nil {
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}
}

func (d *DecisionController) DeleteDecisionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	decisionID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}
	err = d.decisionService.DeleteDecision(uint(decisionID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}
}
