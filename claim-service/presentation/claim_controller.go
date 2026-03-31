package presentation

import (
	"encoding/json"
	"errors"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/janicaleksander/cloud/claimservice/application"
	"github.com/janicaleksander/cloud/claimservice/domain"
)

type ClaimController struct {
	claimService *application.ClaimService
}

func NewClaimController(claimService *application.ClaimService) *ClaimController {
	return &ClaimController{
		claimService: claimService,
	}
}

func parseFile(r *http.Request) ([]*domain.File, error) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return nil, err
	}

	files := r.MultipartForm.File["claim_files"]
	if len(files) == 0 {
		return nil, errors.New("no files provided")
	}

	domainFiles := make([]*domain.File, 0, len(files))
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}

		domainFiles = append(domainFiles, &domain.File{
			FileName: fileHeader.Filename,
			FileExt:  filepath.Ext(fileHeader.Filename),
		})

		file.Close() // close immediately
	}
	return domainFiles, nil
}
func (c *ClaimController) CreateClaimHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	dataField := r.FormValue("json_req_body")
	if dataField == "" {
		http.Error(w, "Missing data field", http.StatusBadRequest)
		return
	}

	var createClaimRequest CreateClaimRequestDTO
	err = json.NewDecoder(strings.NewReader(dataField)).Decode(&createClaimRequest)
	if err != nil {
		http.Error(w, "Invalid JSON in data field", http.StatusBadRequest)
		return
	}

	domainFiles, err := parseFile(r)
	if err != nil {
		http.Error(w, "parsing file error", http.StatusInternalServerError)
		return
	}

	claimDomain := CreateClaimRequestToDomain(&createClaimRequest)
	for idx := range domainFiles {
		domainFiles[idx].StorageURL = "https://storage.example.com/" + domainFiles[idx].FileName
	}
	if len(domainFiles) != 0 {
		claimDomain.Files = domainFiles
	}

	err = c.claimService.CreateClaim(claimDomain)
	if err != nil {
		http.Error(w, "saving Error", http.StatusInternalServerError)
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
