package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	NEW      Status = "NEW"
	VERIFIED Status = "VERIFIED"
	DENIED   Status = "DENIED"
	APPROVED Status = "APPROVED"
	REJECTED Status = "REJECTED"
)

func StringToStatus(s string) (Status, error) {
	switch s {
	case "NEW":
		return NEW, nil
	case "VERIFIED":
		return VERIFIED, nil
	case "DENIED":
		return DENIED, nil
	case "APPROVED":
		return APPROVED, nil
	case "REJECTED":
		return REJECTED, nil
	default:
		return "", errors.New("invalid status")
	}
}

type Claim struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	Email        string
	VIN          string
	AccidentDate time.Time
	Status       Status
	Files        []*File
	UpdatedAt    time.Time
}

type File struct {
	ID         uuid.UUID
	FileName   string
	FileExt    string
	FileSize   int64
	UploadedAt time.Time
	StorageURL string
}

func NewClaim(id, userID uuid.UUID, email, vin string, accidentDate time.Time, status Status, files []*File) *Claim {
	return &Claim{
		ID:           id,
		UserID:       userID,
		Email:        email,
		VIN:          vin,
		AccidentDate: accidentDate,
		Status:       status,
		Files:        files,
		UpdatedAt:    time.Now(),
	}
}

func NewFile(id uuid.UUID, fileName, fileExt string, fileSize int64, uploadedAt time.Time, storageURL string) *File {
	return &File{
		ID:         id,
		FileName:   fileName,
		FileExt:    fileExt,
		FileSize:   fileSize,
		UploadedAt: time.Now(),
		StorageURL: storageURL,
	}
}

type ClaimRepository interface {
	GetAll(context.Context) ([]*Claim, error)
	GetById(context.Context, uuid.UUID) (*Claim, error)
	Save(context.Context, *Claim) (*Claim, error)
	Update(context.Context, *Claim) (*Claim, error)
	UpdateStatus(context.Context, uuid.UUID, Status) error
	DeleteById(context.Context, uuid.UUID) error
	GetFileById(ctx context.Context, fileID uuid.UUID) (*File, error)
}
