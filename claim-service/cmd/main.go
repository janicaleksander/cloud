package main

import (
	"log/slog"

	"github.com/janicaleksander/cloud/claimservice/infrastructure"
	"github.com/janicaleksander/cloud/claimservice/persistance"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file", "error", err)
		panic(err)
	}
	db, err := infrastructure.NewDB()
	if err != nil {
		slog.Error("Error connecting to database", "error", err)
		panic(err)
	}
	err = db.AutoMigrate(&persistance.ClaimModel{}, &persistance.FileModel{})
	if err != nil {
		slog.Error("Error migrating database", "error", err)
		panic(err)
	}

}
