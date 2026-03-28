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
func (c *ClaimController) CreateClaimHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
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

func (c *ClaimController) GetClaimHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
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

	claimDTO := GetClaimDomainToResponse(claim)

	err = json.NewEncoder(w).Encode(&claimDTO)
	if err != nil {
		http.Error(w, "Internal", 500)
		return
	}
}

func (c *ClaimController) GetClaimsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	claims, err := c.claimService.GetClaims()
	if err != nil {
		http.Error(w, "no such claim", 404)
		return
	}
	claimsDTO := make([]*GetClaimResponseDTO, 0, len(claims))
	for idx := range claims {
		claimDTO := GetClaimDomainToResponse(claims[idx])
		claimsDTO = append(claimsDTO, claimDTO)
	}

	err = json.NewEncoder(w).Encode(claimsDTO)
	if err != nil {
		http.Error(w, "Internal", 500)
		return
	}
}

func (c *ClaimController) DeleteClaimHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	claimID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Internal", 500)
		return
	}
	err = c.claimService.DeleteClaim(uint(claimID))
	if err != nil {
		http.Error(w, "no such claim", 404)
		return
	}

	err = json.NewEncoder(w).Encode(map[string]string{"message": "claim deleted" + strconv.Itoa(claimID)})
	if err != nil {
		http.Error(w, "Internal", 500)
		return
	}
}

func (c *ClaimController) UpdateClaimHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	claimID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Internal", 500)
		return
	}
	var updateClaimRequestDTO UpdateClaimRequestDTO
	err = json.NewDecoder(r.Body).Decode(&updateClaimRequestDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}

	claim, err := c.claimService.GetClaim(uint(claimID))
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	if updateClaimRequestDTO.UserID != claim.UserID && updateClaimRequestDTO.UserID != 0 {
		claim.UserID = updateClaimRequestDTO.UserID
	}
	if updateClaimRequestDTO.CarID != claim.CarID && updateClaimRequestDTO.CarID != 0 {
		claim.CarID = updateClaimRequestDTO.CarID
	}
	err = c.claimService.UpdateClaim(claim)
	if err != nil {
		http.Error(w, err.Error(), 403)
		return
	}
	claimResponse := GetClaimDomainToResponse(claim)
	err = json.NewEncoder(w).Encode(&claimResponse)
	if err != nil {
		http.Error(w, err.Error(), 403)
		return
	}
}
