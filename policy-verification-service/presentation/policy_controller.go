package presentation

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/janicaleksander/cloud/policyverificationservice/application"
)

type PolicyController struct {
	policyService *application.PolicyService
}

func NewPolicyController(policyService *application.PolicyService) *PolicyController {
	return &PolicyController{policyService: policyService}
}

// TODO throw event when i check if user has policy
func (p *PolicyController) CreatePolicyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var d CreatePolicyRequestDTO
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}
	policyDomain := CreatePolicyRequestToDomain(&d)
	err = p.policyService.CreatePolicy(policyDomain)
	if err != nil {
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}

}

func (p *PolicyController) GetPolicyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	policyId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid policy ID", http.StatusBadRequest)
		return

	}
	domainPolicy, err := p.policyService.GetPolicy(uint(policyId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	d := GetPolicyDomainToResponse(domainPolicy)
	err = json.NewEncoder(w).Encode(d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (p *PolicyController) GetPoliciesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	domainPolicies, err := p.policyService.GetPolicies()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	responseDTOs := make([]*GetPolicyResponseDTO, 0, len(domainPolicies))
	for idx := range domainPolicies {
		responseDTOs = append(responseDTOs, GetPolicyDomainToResponse(domainPolicies[idx]))
	}
	err = json.NewEncoder(w).Encode(responseDTOs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (p *PolicyController) UpdatePolicyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	policyId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid policy ID", http.StatusBadRequest)
		return
	}
	var d UpdatePolicyRequest
	err = json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	policyDomain, err := p.policyService.GetPolicy(uint(policyId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = p.policyService.UpdatePolicy(policyDomain, d.From, d.To)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (p *PolicyController) DeletePolicyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	policyId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid policy ID", http.StatusBadRequest)
		return
	}
	err = p.policyService.DeletePolicy(uint(policyId))
	if err != nil {
		http.Error(w, "no such policy ID", http.StatusBadRequest)
		return
	}
}
