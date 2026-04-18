package interfaces

import (
	"context"
	"io"
)

type FileStorage interface {
	StoreFile(ctx context.Context, bucket string, fileID string, contentType string, reader io.Reader) error
	GetFile(ctx context.Context, bucket string, fileID string) (io.ReadCloser, error)
	RemoveFile(ctx context.Context, bucket string, fileID string) error
}

type ClaimEventPublisher interface {
	Publish(exchange string, msg interface{}) error
}
