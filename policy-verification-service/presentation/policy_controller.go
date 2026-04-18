package presentation

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/policyverificationservice/application/command"
	"github.com/janicaleksander/cloud/policyverificationservice/application/query"
	"github.com/mehdihadeli/go-mediatr"
)

type PolicyController struct {
}

func NewPolicyController() *PolicyController {
	slog.Info("Creating PolicyController")
	return &PolicyController{}
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
	uid := uuid.New()
	cmd := CreatePolicyResponseHTTPToCommand(uid, &d)

	_, err = mediatr.Send[*command.CreatePolicyCommand, *mediatr.Unit](context.Background(), cmd)
	if err != nil {
		failure(w, http.StatusInternalServerError, "Error creating policy: "+err.Error())
		return
	}
	successWithLocation(w, "/policy/"+uid.String(), http.StatusCreated)

}

// TODO repiar in dtos that i have embeded {} with the same tag in every get
func (p *PolicyController) GetPolicyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		failure(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	slog.Info("HTTP GetPolicyHandler called")
	policyId := chi.URLParam(r, "id")
	q := &query.GetPolicyQuery{PolicyID: policyId}

	d, err := mediatr.Send[*query.GetPolicyQuery, *query.GetPolicyQueryResponse](context.Background(), q)
	if err != nil {
		failure(w, http.StatusInternalServerError, "Error fetching policy: "+err.Error())
		return
	}
	success(w, map[string]any{"policy": d}, http.StatusOK)

}

func (p *PolicyController) GetPoliciesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		failure(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	slog.Info("HTTP GetPoliciesHandler called")
	q := &query.GetPoliciesQuery{}
	policies, err := mediatr.Send[*query.GetPoliciesQuery, *query.GetPoliciesQueryResponse](context.Background(), q)
	if err != nil {
		failure(w, http.StatusInternalServerError, "Error fetching policies: "+err.Error())
		return
	}
	success(w, map[string]any{"policies": policies}, http.StatusOK)
}

func (p *PolicyController) UpdatePolicyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		failure(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	slog.Info("HTTP UpdatePolicyHandler called")
	policyId := chi.URLParam(r, "id")

	var d UpdatePolicyRequest
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		failure(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	cmd := command.UpdatePolicyCommand{
		PolicyID: policyId,
		NewFrom:  d.From,
		NewTo:    d.From,
	}
	_, err = mediatr.Send[*command.UpdatePolicyCommand, *mediatr.Unit](context.Background(), &cmd)
	if err != nil {
		failure(w, http.StatusInternalServerError, "Error updating policy: "+err.Error())
		return
	}
	successWithLocation(w, "/policy/"+policyId, http.StatusOK)
}

func (p *PolicyController) DeletePolicyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		failure(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	slog.Info("HTTP DeletePolicyHandler called")
	policyId := chi.URLParam(r, "id")

	cmd := command.DeletePolicyCommand{PolicyID: policyId}
	_, err := mediatr.Send[*command.DeletePolicyCommand, *mediatr.Unit](context.Background(), &cmd)
	if err != nil {
		failure(w, http.StatusInternalServerError, "Error deleting policy: "+err.Error())
		return
	}
	success(w, nil, http.StatusNoContent)
}
