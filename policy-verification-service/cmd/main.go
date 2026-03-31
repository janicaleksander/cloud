package main

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/janicaleksander/cloud/common/rabbitmq"
	"github.com/janicaleksander/cloud/policyverificationservice/application"
	"github.com/janicaleksander/cloud/policyverificationservice/infrastructure"
	"github.com/janicaleksander/cloud/policyverificationservice/infrastructure/messaging"
	"github.com/janicaleksander/cloud/policyverificationservice/persistance"
	"github.com/janicaleksander/cloud/policyverificationservice/presentation/api"
	"github.com/janicaleksander/cloud/policyverificationservice/presentation/api/router"
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
	err = db.AutoMigrate(&persistance.PolicyModel{})
	if err != nil {
		panic(err)
	}
	rabbit, err := rabbitmq.NewRabbitMQ()
	if err != nil {
		panic(err)
	}
	publisher := rabbitmq.NewPublisher(rabbit)

	policyRepository := persistance.NewPolicyRepository(db)
	policyService := application.NewPolicyService(policyRepository, publisher)
	policyController := api.NewPolicyController(policyService)

	policyEventHandler := messaging.NewPolicyEventHandler(policyService)
	policyEventHandler.Run(rabbit)
	//evsent handler

	r := router.NewRouter(policyController)
	log.Println("serving on 8081")
	err = http.ListenAndServe("localhost:8081", r)
	if err != nil {
		slog.Error("Error running http: ", err)
		panic(err)
	}
}
