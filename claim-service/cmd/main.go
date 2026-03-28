package main

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/janicaleksander/cloud/claimservice/application"
	"github.com/janicaleksander/cloud/claimservice/infrastructure"
	"github.com/janicaleksander/cloud/claimservice/persistance"
	"github.com/janicaleksander/cloud/claimservice/presentation"
	"github.com/janicaleksander/cloud/claimservice/presentation/api/router"
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
	claimService := application.NewClaimService(
		persistance.NewClaimRepository(db),
	)
	claimController := presentation.NewClaimController(claimService)
	r := router.NewRouter(claimController)
	log.Println("serving on 8080")
	err = http.ListenAndServe("localhost:8080", r)
	if err != nil {
		slog.Error("Error running http: ", err)
		panic(err)
	}

}
