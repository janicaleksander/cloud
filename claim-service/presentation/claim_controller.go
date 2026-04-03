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
		failure(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		failure(w, http.StatusBadRequest, "Invalid multipart form data")
		return
	}

	dataField := r.FormValue("json_req_body")
	if dataField == "" {
		failure(w, http.StatusBadRequest, "Missing json_req_body field")
		return
	}

	var createClaimRequest CreateClaimRequestDTO
	err = json.NewDecoder(strings.NewReader(dataField)).Decode(&createClaimRequest)
	if err != nil {
		failure(w, http.StatusBadRequest, "Invalid JSON in json_req_body")
		return
	}

	domainFiles, err := parseFile(r)
	if err != nil {
		failure(w, http.StatusBadRequest, "Error processing files: "+err.Error())
		return
	}

	claimDomain := CreateClaimRequestToDomain(&createClaimRequest)
	for idx := range domainFiles {
		domainFiles[idx].StorageURL = "https://storage.example.com/" + domainFiles[idx].FileName
	}
	if len(domainFiles) != 0 {
		claimDomain.Files = domainFiles
	}

	createdClaim, err := c.claimService.CreateClaim(claimDomain)
	if err != nil {
		failure(w, http.StatusInternalServerError, "Error creating claim: "+err.Error())
		return
	}
	dto := GetClaimDomainToResponse(createdClaim)
	success(w, dto)
}

func (c *ClaimController) GetClaimHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		failure(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	claimID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		failure(w, http.StatusBadRequest, "Invalid claim ID")
		return
	}
	claim, err := c.claimService.GetClaim(uint(claimID))
	if err != nil {
		failure(w, http.StatusNotFound, "No such claim")
		return
	}

	claimDTO := GetClaimDomainToResponse(claim)

	success(w, map[string]any{"claim": claimDTO})
}

func (c *ClaimController) GetClaimsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		failure(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	claims, err := c.claimService.GetClaims()
	if err != nil {
		failure(w, http.StatusInternalServerError, "Error fetching claims: "+err.Error())
		return
	}
	claimsDTO := make([]*GetClaimResponseDTO, 0, len(claims))
	for idx := range claims {
		claimDTO := GetClaimDomainToResponse(claims[idx])
		claimsDTO = append(claimsDTO, claimDTO)
	}

	success(w, map[string]any{"claims": claimsDTO})
}

func (c *ClaimController) DeleteClaimHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		failure(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	claimID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		failure(w, http.StatusBadRequest, "Invalid claim ID")
		return
	}
	err = c.claimService.DeleteClaim(uint(claimID))
	if err != nil {
		failure(w, http.StatusInternalServerError, "Error deleting claim: "+err.Error())
		return
	}
	success(w, map[string]any{"message": "Claim deleted successfully +: " + strconv.Itoa(claimID)})
}

func (c *ClaimController) UpdateClaimHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		failure(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	claimID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		failure(w, http.StatusBadRequest, "Invalid claim ID")
		return
	}
	var updateClaimRequestDTO UpdateClaimRequestDTO
	err = json.NewDecoder(r.Body).Decode(&updateClaimRequestDTO)
	if err != nil {
		failure(w, http.StatusBadRequest, "Invalid JSON in request body")
		return
	}

	claim, err := c.claimService.GetClaim(uint(claimID))
	if err != nil {
		failure(w, http.StatusNotFound, "No such claim")
		return
	}

	updatedClaim, err := c.claimService.UpdateClaim(claim, updateClaimRequestDTO.Email)
	if err != nil {
		failure(w, http.StatusForbidden, "Error updating claim: "+err.Error())
		return
	}
	claimResponse := GetClaimDomainToResponse(updatedClaim)
	success(w, map[string]any{"claim": claimResponse})
}
