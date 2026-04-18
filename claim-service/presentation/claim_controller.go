package presentation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/claimservice/application/command"
	"github.com/janicaleksander/cloud/claimservice/application/query"
	"github.com/mehdihadeli/go-mediatr"
)

type ClaimController struct {
}

func NewClaimController() *ClaimController {
	slog.Info("Creating ClaimController")
	return &ClaimController{}
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

func parseFile(r *http.Request) ([]*os.File, error) {
	slog.Info("Parsing multipart form data for files")
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return nil, err
	}

	files := r.MultipartForm.File["claim_files"]
	if len(files) == 0 {
		return nil, errors.New("no files provided")
	}

	filesObject := make([]*os.File, 0, len(files))
	for _, fileHeader := range files {
		multipartFile, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}

		slog.Info("Processing file: ", "filename", fileHeader.Filename, "size", fileHeader.Size)
		f, err := os.Create(fileHeader.Filename)
		if err != nil {
			multipartFile.Close()
			return nil, err
		}

		// Copy contents from the uploaded file to the local temp file
		_, err = io.Copy(f, multipartFile)
		multipartFile.Close() // close the multipart file after copying
		if err != nil {
			f.Close()
			return nil, err
		}

		// Reset the file pointer to the beginning so it can be read later by S3
		_, err = f.Seek(0, 0)
		if err != nil {
			f.Close()
			return nil, err
		}

		filesObject = append(filesObject, f)
	}
	return filesObject, nil
}
func (c *ClaimController) CreateClaimHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		failure(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	slog.Info("HTTP CreateClaimHandler called")

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

	objectFiles, err := parseFile(r)
	if err != nil {
		failure(w, http.StatusBadRequest, "Error processing files: "+err.Error())
		return
	}

	claimDomain := HTTPCreateClaimRequestToDomain(&createClaimRequest)
	claimDomain.ID = uuid.New()
	cmd := command.ClaimDomainToCreateClaimCommand(claimDomain, objectFiles)

	_, err = mediatr.Send[*command.CreateClaimCommand, *mediatr.Unit](context.Background(), cmd)
	if err != nil {
		failure(w, http.StatusInternalServerError, "Error creating claim: "+err.Error())
		return
	}
	successWithLocation(w, fmt.Sprintf("/claim/%s", claimDomain.ID), http.StatusCreated)

}

func (c *ClaimController) GetClaimHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		failure(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	slog.Info("HTTP GetClaimHandler called")
	claimID := chi.URLParam(r, "id")
	q := &query.GetClaimByIdQuery{ClaimID: claimID}
	response, err := mediatr.Send[*query.GetClaimByIdQuery, *query.GetClaimByIdQueryResponse](context.Background(), q)
	if err != nil {
		failure(w, http.StatusNotFound, "No such claim")
		return
	}
	claimDomain := query.GetClaimQueryResponseToDomain(response)

	claimDTO := HTTPGetClaimDomainToResponse(claimDomain)
	success(w, map[string]any{"claim": claimDTO}, 200)
}

func (c *ClaimController) GetClaimsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		failure(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	slog.Info("HTTP GetClaimsHandler called")
	q := &query.GetClaimsQuery{}
	response, err := mediatr.Send[*query.GetClaimsQuery, *query.GetClaimsQueryResponse](context.Background(), q)
	if err != nil {
		failure(w, http.StatusInternalServerError, "Error fetching claims: "+err.Error())
		return
	}

	claims := query.GetClaimsQueryResponseToDomain(response.Claims)

	claimsDTO := make([]*GetClaimResponseDTO, 0, len(claims))
	for idx := range claims {
		claimDTO := HTTPGetClaimDomainToResponse(claims[idx])
		claimsDTO = append(claimsDTO, claimDTO)
	}

	success(w, map[string]any{"claims": claimsDTO}, 200)
}

func (c *ClaimController) DeleteClaimHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		failure(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	slog.Info("HTTP DeleteClaimHandler called")
	claimID := chi.URLParam(r, "id")
	cmd := &command.DeleteClaimCommand{ClaimID: claimID}
	_, err := mediatr.Send[*command.DeleteClaimCommand, *mediatr.Unit](context.Background(), cmd)
	if err != nil {
		failure(w, http.StatusInternalServerError, "Error deleting claim: "+err.Error())
		return
	}
	successWithLocation(w, fmt.Sprintf("/claim/%s", claimID), http.StatusNoContent)

}

func (c *ClaimController) UpdateClaimHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		failure(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	slog.Info("HTTP UpdateClaimHandler called")
	claimID := chi.URLParam(r, "id")
	var updateClaimRequestDTO UpdateClaimRequestDTO
	err := json.NewDecoder(r.Body).Decode(&updateClaimRequestDTO)
	if err != nil {
		failure(w, http.StatusBadRequest, "Invalid JSON in request body")
		return
	}
	cmd := &command.UpdateClaimCommand{
		ClaimID:  claimID,
		NewEmail: updateClaimRequestDTO.Email,
	}
	_, err = mediatr.Send[*command.UpdateClaimCommand, *mediatr.Unit](context.Background(), cmd)
	if err != nil {
		failure(w, http.StatusNotFound, "No such claim")
		return
	}
	successWithLocation(w, fmt.Sprintf("/claim/%s", claimID), http.StatusOK)

}

func (c *ClaimController) GetFileFromStorageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		failure(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	fileID := chi.URLParam(r, "id")
	q := &query.GetFileFromStorageQuery{FileID: fileID}
	response, err := mediatr.Send[*query.GetFileFromStorageQuery, *query.GetFileFromStorageQueryResponse](context.Background(), q)
	if err != nil {
		failure(w, http.StatusNotFound, "No such file")
		return
	}

	contentType := mime.TypeByExtension(response.FileExt)
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", response.FileName))

	defer response.Reader.Close()
	io.Copy(w, response.Reader)
}
