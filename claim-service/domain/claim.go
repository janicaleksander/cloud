package domain

import (
	"gorm.io/gorm"
)

type Status string

const (
	NEW      Status = "NEW"
	VERIFIED Status = "VERIFIED"
	DENIED   Status = "DENIED"
	APPROVED Status = "APPROVED"
	REJECTED Status = "REJECTED"
)

type Claim struct {
	gorm.Model
	UserID int    `gorm:"not null"`
	CarID  int    `gorm:"not null"`
	Status string `gorm:"not null"`
	Files  []File
}

type File struct {
	gorm.Model
	ClaimID  int    `gorm:"not null"`
	FileName string `gorm:"not null"`
	FileExt  string `gorm:"not null"`
}

//evnet driven architecture ? hanlder" a nie w ramach serwisu

//jakis common serwis ktory ma nprabbit mq config  i on jest zaciagany  przez inne mikrosweriy
//ten event -> ten handler
