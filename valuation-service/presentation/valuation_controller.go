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

func (v *ValuationController) GetValuationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		failure(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	var d GetValuationResponseDTO

	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		failure(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	domainValuation, err := v.valuationService.GetValuations()
	if err != nil {
		failure(w, http.StatusInternalServerError, "Failed to get valuations")
		return
	}
	valuationDTO := make([]*GetValuationResponseDTO, 0, len(domainValuation))
	for _, valu := range domainValuation {
		valuationDTO = append(valuationDTO, GetValuationDomainToResponse(valu))
	}
	success(w, map[string]any{"valuations": valuationDTO})

}

func (v *ValuationController) GetValuationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		failure(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	var d GetValuationResponseDTO

	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		failure(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	domainValuation, err := v.valuationService.GetValuation(d.ClaimID)
	if err != nil {
		failure(w, http.StatusInternalServerError, "Failed to get valuation")
		return
	}
	valuationDTO := GetValuationDomainToResponse(domainValuation)
	success(w, map[string]any{"valuation": valuationDTO})
}

/*
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
*/
func (v *ValuationController) DeleteValuationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		failure(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	valuationID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		failure(w, http.StatusBadRequest, "Invalid valuation ID")
		return
	}
	err = v.valuationService.DeleteValuation(uint(valuationID))
	if err != nil {
		failure(w, http.StatusInternalServerError, "Failed to delete valuation")
		return
	}
	success(w, map[string]any{"message": "Valuation deleted successfully"})

}
