package domain

import (
	"context"
	"time"
)

type MetaFile struct {
	ID       string
	FileName string
	FileExt  string
	FileSize float64
	Date     time.Time
	FileURL  string
}

type MetaFileRepository interface {
	Create(ctx context.Context, file *MetaFile) (*MetaFile, error)
	GetFileById(ctx context.Context, id string) (*MetaFile, error)
	GetFiles(ctx context.Context) ([]*MetaFile, error)
	DeleteFileById(ctx context.Context, id string) error
}
