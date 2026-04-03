package main

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/janicaleksander/cloud/common/rabbitmq"
	"github.com/janicaleksander/cloud/decisionservice/application"
	"github.com/janicaleksander/cloud/decisionservice/infrastructure"
	"github.com/janicaleksander/cloud/decisionservice/infrastructure/messaging"
	"github.com/janicaleksander/cloud/decisionservice/persistence"
	"github.com/janicaleksander/cloud/decisionservice/presentation"
	"github.com/janicaleksander/cloud/decisionservice/presentation/router"
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
	err = db.AutoMigrate(&persistence.DecisionModel{})
	if err != nil {
		panic(err)
	}

	rabbit, err := rabbitmq.NewRabbitMQ()
	if err != nil {
		panic(err)
	}
	publisher := rabbitmq.NewPublisher(rabbit)
	decisionRepository := persistence.NewDecisionRepository(db)
	decisionService := application.NewDecisionService(decisionRepository, publisher)
	decisionController := presentation.NewDecisionController(decisionService)
	decisionEventHandler := messaging.NewDecisionEventHandler(decisionService)
	err = decisionEventHandler.Run(rabbit)
	if err != nil {
		panic(err)
	}

	r := router.NewRouter(decisionController)
	log.Println("serving on 8084")
	err = http.ListenAndServe("localhost:8084", r)
	if err != nil {
		slog.Error("Error running http: ", err)
		panic(err)
	}

}
