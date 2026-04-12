package presentation

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/janicaleksander/cloud/policyverificationservice/application"
)

type PolicyController struct {
	policyService *application.PolicyService
}

func NewPolicyController(policyService *application.PolicyService) *PolicyController {
	slog.Info("Creating PolicyController")
	return &PolicyController{policyService: policyService}
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

func (p *PolicyController) CreatePolicyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		failure(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	slog.Info("HTTP CreatePolicyHandler called")
	var d CreatePolicyRequestDTO
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		failure(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	policyDomain := CreatePolicyRequestToDomain(&d)
	createdPolicy, err := p.policyService.CreatePolicy(policyDomain)
	if err != nil {
		failure(w, http.StatusInternalServerError, err.Error())
		return
	}
	success(w, map[string]any{"policy": GetPolicyDomainToResponse(createdPolicy)})

}

func (p *PolicyController) GetPolicyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		failure(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	slog.Info("HTTP GetPolicyHandler called")
	policyId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		failure(w, http.StatusBadRequest, "Invalid policy ID")
		return

	}
	domainPolicy, err := p.policyService.GetPolicy(uint(policyId))
	if err != nil {
		failure(w, http.StatusBadRequest, "no such policy ID")
		return
	}
	d := GetPolicyDomainToResponse(domainPolicy)
	success(w, map[string]any{"policy": d})

}

func (p *PolicyController) GetPoliciesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		failure(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	slog.Info("HTTP GetPoliciesHandler called")
	domainPolicies, err := p.policyService.GetPolicies()
	if err != nil {
		failure(w, http.StatusInternalServerError, "Error fetching policies: "+err.Error())
		return
	}
	responseDTOs := make([]*GetPolicyResponseDTO, 0, len(domainPolicies))
	for idx := range domainPolicies {
		responseDTOs = append(responseDTOs, GetPolicyDomainToResponse(domainPolicies[idx]))
	}
	success(w, map[string]any{"policies": responseDTOs})
}

func (p *PolicyController) UpdatePolicyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		failure(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	slog.Info("HTTP UpdatePolicyHandler called")
	policyId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		failure(w, http.StatusBadRequest, "Invalid policy ID")
		return
	}
	var d UpdatePolicyRequest
	err = json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		failure(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	policyDomain, err := p.policyService.GetPolicy(uint(policyId))
	if err != nil {
		failure(w, http.StatusBadRequest, "no such policy ID")
		return
	}

	updatedPolicy, err := p.policyService.UpdatePolicy(policyDomain, d.From, d.To)
	if err != nil {
		failure(w, http.StatusInternalServerError, "Error updating policy: "+err.Error())
		return
	}
	success(w, map[string]any{"policy": GetPolicyDomainToResponse(updatedPolicy)})
}

func (p *PolicyController) DeletePolicyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		failure(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	slog.Info("HTTP DeletePolicyHandler called")
	policyId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		failure(w, http.StatusBadRequest, "Invalid policy ID")
		return
	}
	err = p.policyService.DeletePolicy(uint(policyId))
	if err != nil {
		failure(w, http.StatusInternalServerError, "Error deleting policy: "+err.Error())
		return
	}
	success(w, map[string]any{"message": "policy deleted" + strconv.Itoa(policyId)})
}
