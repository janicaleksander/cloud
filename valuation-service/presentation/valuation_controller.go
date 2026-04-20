package presentation

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/janicaleksander/cloud/valuationservice/application/command"
	"github.com/janicaleksander/cloud/valuationservice/application/query"
	"github.com/mehdihadeli/go-mediatr"
)

type ValuationController struct {
}

func NewValuationController() *ValuationController {
	slog.Info("Creating ValuationController")
	return &ValuationController{}
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

func (v *ValuationController) GetValuationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		failure(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	slog.Info("HTTP GetValuationsHandler called")
	q := &query.GetValuationsQuery{}
	valuationsResponse, err := mediatr.Send[*query.GetValuationsQuery, *query.GetValuationsQueryResponse](context.Background(), q)
	if err != nil {
		failure(w, http.StatusInternalServerError, "Failed to get valuations")
		return
	}
	success(w, map[string]any{"valuations": valuationsResponse.Valuations}, 200)
}

func (v *ValuationController) GetValuationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		failure(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	slog.Info("HTTP GetValuationHandler called")
	valuationID := chi.URLParam(r, "id")
	q := &query.GetValuationQuery{ClaimID: valuationID}
	valuationResponse, err := mediatr.Send[*query.GetValuationQuery, *query.GetValuationQueryResponse](context.Background(), q)
	if err != nil {
		failure(w, http.StatusInternalServerError, "Failed to get valuation")
		return
	}
	success(w, map[string]any{"valuation": valuationResponse}, 200)
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
	slog.Info("HTTP DeleteValuationHandler called")
	valuationID := chi.URLParam(r, "id")

	cmd := &command.DeleteValuationCommand{ID: valuationID}
	_, err := mediatr.Send[*command.DeleteValuationCommand, *mediatr.Unit](context.Background(), cmd)
	if err != nil {
		failure(w, http.StatusInternalServerError, "Failed to delete valuation")
		return
	}
	success(w, map[string]any{"message": "Valuation deleted successfully"}, 204)
}
