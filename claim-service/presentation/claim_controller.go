package presentation

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/janicaleksander/cloud/claimservice/application"
)

type ClaimController struct {
	claimService *application.ClaimService
}

func NewClaimController(claimService *application.ClaimService) *ClaimController {
	return &ClaimController{
		claimService: claimService,
	}
}

func (c *ClaimController) GetClaimHandler(w http.ResponseWriter, r *http.Request) {
	claimID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Internal", 500)
		return
	}
	claim, err := c.claimService.GetClaim(uint(claimID))
	if err != nil {
		http.Error(w, "no such claim", 404)
		return
	}

	claimDTO := GetClaimDomainToRequest(claim)

	err = json.NewEncoder(w).Encode(&claimDTO)
	if err != nil {
		http.Error(w, "Internal", 500)
		return
	}
}

func (c *ClaimController) GetClaimsHandler(w http.ResponseWriter, r *http.Request) {

}

func (c *ClaimController) CreateClaimHandler(w http.ResponseWriter, r *http.Request) {
	var createClaimRequest CreateClaimRequestDTO
	err := json.NewDecoder(r.Body).Decode(&createClaimRequest)
	if err != nil {
		http.Error(w, "Error", 408)
		return
	}

	claimDomain := CreateClaimRequestToDomain(&createClaimRequest)

	err = c.claimService.CreateClaim(claimDomain)
	if err != nil {
		http.Error(w, "saving Error", 408)
		return
	}

	err = json.NewEncoder(w).Encode(createClaimRequest)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
	}

}
