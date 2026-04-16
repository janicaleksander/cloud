package domain

import (
	"context"
	"errors"
	"time"
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
	ID           uint
	UserID       uint
	Email        string
	VIN          string
	AccidentDate time.Time
	Status       Status
	Files        []*File
	UpdatedAt    time.Time
}

type File struct {
	ID         uint
	FileName   string
	FileExt    string
	FileSize   int64
	UploadedAt time.Time
	StorageURL string
}

type ClaimRepository interface {
	GetAll(context.Context) ([]*Claim, error)
	GetById(context.Context, uint) (*Claim, error)
	Save(context.Context, *Claim) (*Claim, error)
	Update(context.Context, *Claim) (*Claim, error)
	UpdateStatus(context.Context, uint, Status) error
	DeleteById(context.Context, uint) error
	GetFileById(ctx context.Context, fileID uint) (*File, error)
}
