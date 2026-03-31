package main

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/janicaleksander/cloud/common/rabbitmq"
	"github.com/janicaleksander/cloud/valuationservice/application"
	"github.com/janicaleksander/cloud/valuationservice/infrastructure"
	"github.com/janicaleksander/cloud/valuationservice/infrastructure/messaging"
	"github.com/janicaleksander/cloud/valuationservice/persistance"
	presentation "github.com/janicaleksander/cloud/valuationservice/presentation/api"
	"github.com/janicaleksander/cloud/valuationservice/presentation/api/router"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	db, err := infrastructure.NewDB()
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&persistance.ValuationModel{})
	if err != nil {
		panic(err)
	}

	rabbit, err := rabbitmq.NewRabbitMQ()
	publisher := rabbitmq.NewPublisher(rabbit)

	if err != nil {
		panic(err)
	}
	valuationRepository := persistance.NewValuationRepository(db)
	valuationService := application.NewValuationService(valuationRepository, publisher)
	valuationController := presentation.NewValuationController(valuationService)
	valuationHandler := messaging.NewValuationEventHandler(valuationService)
	valuationHandler.Run(rabbit)
	r := router.NewRouter(valuationController)

	log.Println("serving on 8082")
	err = http.ListenAndServe("localhost:8082", r)
	if err != nil {
		slog.Error("Error running http: ", err)
		panic(err)
	}

}
