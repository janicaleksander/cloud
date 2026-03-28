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
	ID        uint
	UserID    uint
	CarID     uint
	Status    Status
	Files     []*File
	CreatedAt time.Time
	UpdatedAt time.Time
}

type File struct {
	ID         uint
	FileName   string
	FileExt    string
	StorageURL string //added
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type ClaimerRepository interface {
	GetAll(context.Context) ([]*Claim, error)
	GetById(context.Context, uint) (*Claim, error)
	Save(context.Context, *Claim) error
	Update(context.Context, *Claim) error
	DeleteById(context.Context, uint) error
}

//evnet driven architecture ? hanlder" a nie w ramach serwisu

//jakis common serwis ktory ma nprabbit mq config  i on jest zaciagany  przez inne mikrosweriy
//ten event -> ten handler
