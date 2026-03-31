package presentation

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/janicaleksander/cloud/valuationservice/application"
)

type ValuationController struct {
	valuationService *application.ValuationService
}

func NewValuationController(vS *application.ValuationService) *ValuationController {
	return &ValuationController{
		valuationService: vS,
	}
}
func (v *ValuationController) GetValuationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var d GetValuationResponseDTO

	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	domainValuation, err := v.valuationService.GetValuations()
	if err != nil {
		http.Error(w, "Failed to get valuations", http.StatusInternalServerError)
		return
	}
	valuationDTO := make([]*GetValuationResponseDTO, 0, len(domainValuation))
	for _, valu := range domainValuation {
		valuationDTO = append(valuationDTO, GetValuationDomainToResponse(valu))
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(valuationDTO)

}

func (v *ValuationController) GetValuationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var d GetValuationResponseDTO

	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	domainValuation, err := v.valuationService.GetValuation(d.ClaimID)
	if err != nil {
		http.Error(w, "Failed to get valuation", http.StatusInternalServerError)
		return
	}
	valuationDTO := GetValuationDomainToResponse(domainValuation)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(valuationDTO)
}

func (v *ValuationController) UpdateValuationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	valuationID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid valuation ID", http.StatusBadRequest)
		return
	}
	var d UpdateValuationRequestDTO
	err = json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	oldValuation, err := v.valuationService.GetValuation(uint(valuationID))
	if err != nil {
		http.Error(w, "Failed to get valuation", http.StatusInternalServerError)
		return
	}
	_, err = v.valuationService.UpdateValuation(oldValuation, d.Amount)
	if err != nil {
		http.Error(w, "Failed to update valuation", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}

func (v *ValuationController) DeleteValuationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	valuationID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid valuation ID", http.StatusBadRequest)
		return
	}
	err = v.valuationService.DeleteValuation(uint(valuationID))
	if err != nil {
		http.Error(w, "Failed to delete valuation", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
