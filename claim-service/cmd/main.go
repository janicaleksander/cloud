package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/janicaleksander/cloud/claimservice/application/command"
	"github.com/janicaleksander/cloud/claimservice/application/query"
	"github.com/janicaleksander/cloud/claimservice/infrastructure"
	"github.com/janicaleksander/cloud/claimservice/infrastructure/aws"
	"github.com/janicaleksander/cloud/claimservice/infrastructure/messaging"
	"github.com/janicaleksander/cloud/claimservice/persistence"
	"github.com/janicaleksander/cloud/claimservice/presentation"
	"github.com/janicaleksander/cloud/claimservice/presentation/router"
	"github.com/janicaleksander/cloud/common/rabbitmq"
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

	err = db.AutoMigrate(&persistence.ClaimModel{}, &persistence.FileModel{})
	if err != nil {
		slog.Error("Error migrating database", "error", err)
		panic(err)
	}

	rabbit, err := rabbitmq.NewRabbitMQ()
	if err != nil {
		panic(err)
	}
	publisher := rabbitmq.NewPublisher(rabbit)

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-east-1"))
	if err != nil {
		panic(err)
	}
	client := s3.NewFromConfig(cfg)
	fileStorage := aws.NewAWSStorage(client)

	claimRepo := persistence.NewClaimRepository(db)

	createClaimHandler := command.NewCreateClaimCommandHandler(claimRepo, publisher, fileStorage)
	deleteClaimHandler := command.NewDeleteClaimCommandHandler(claimRepo, fileStorage)
	updateStatusHandler := command.NewUpdateClaimStatusCommandHandler(claimRepo, publisher)

	_ = createClaimHandler.SelfRegister()
	_ = deleteClaimHandler.SelfRegister()
	_ = updateStatusHandler.SelfRegister()

	getClaimQuery := query.NewGetClaimQueryHandler(claimRepo)
	getClaimsQuery := query.NewGetClaimsQueryHandler(claimRepo)
	getFileQuery := query.NewGetFileFromStorageQueryHandler(claimRepo, fileStorage)

	_ = getClaimQuery.SelfRegister()
	_ = getClaimsQuery.SelfRegister()
	_ = getFileQuery.SelfRegister()
	claimEventHandler := messaging.NewClaimEventHandler()
	go claimEventHandler.Run(rabbit)

	// =========================
	// HTTP
	// =========================
	claimController := presentation.NewClaimController()
	r := router.NewRouter(claimController)

	log.Println("serving on 8080")
	err = http.ListenAndServe("localhost:8080", r)
	if err != nil {
		slog.Error("Error running http: ", err)
		panic(err)
	}
}
