package query

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/claimservice/application/interfaces"
	"github.com/janicaleksander/cloud/claimservice/domain"
	"github.com/mehdihadeli/go-mediatr"
)

func ParseS3URL(rawURL string) (bucket, key string, err error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", "", fmt.Errorf("invalid URL: %w", err)
	}

	// Host: "{bucket}.s3.{region}.amazonaws.com"
	host := u.Hostname()
	parts := strings.SplitN(host, ".", 2)
	if len(parts) < 2 || !strings.HasSuffix(parts[1], ".amazonaws.com") {
		return "", "", fmt.Errorf("not an S3 URL: %s", host)
	}

	bucket = parts[0]
	key = strings.TrimPrefix(u.Path, "/")

	return bucket, key, nil
}

type GetFileFromStorageQuery struct {
	FileID string `json:"file_id"`
}

type GetFileFromStorageQueryResponse struct {
	FileExt  string
	FileName string
	Reader   io.ReadCloser
}

type GetFileFromStorageQueryHandler struct {
	repo        domain.ClaimRepository
	fileStorage interfaces.FileStorage
}

func NewGetFileFromStorageQueryHandler(r domain.ClaimRepository) *GetFileFromStorageQueryHandler {
	return &GetFileFromStorageQueryHandler{repo: r}

}

func (h *GetFileFromStorageQueryHandler) SelfRegister() error {
	return mediatr.RegisterRequestHandler[*GetFileFromStorageQuery, *GetFileFromStorageQueryResponse](h)

}

func (h *GetFileFromStorageQueryHandler) Handle(ctx context.Context, query *GetFileFromStorageQuery) (*GetFileFromStorageQueryResponse, error) {
	fid, err := uuid.Parse(query.FileID)
	if err != nil {
		return nil, err
	}
	file, err := h.repo.GetFileById(ctx, fid)
	if err != nil {
		return nil, err
	}
	bucket, key, err := ParseS3URL(file.StorageURL)
	if err != nil {
		return nil, err
	}
	reader, err := h.fileStorage.GetFile(context.Background(), bucket, key)
	if err != nil {
		return nil, err
	}
	return &GetFileFromStorageQueryResponse{
		FileExt:  file.FileExt,
		FileName: file.FileName,
		Reader:   reader,
	}, nil
}
