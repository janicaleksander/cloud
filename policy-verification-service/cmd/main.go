package main

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/janicaleksander/cloud/common/rabbitmq"
	"github.com/janicaleksander/cloud/policyverificationservice/application/command"
	"github.com/janicaleksander/cloud/policyverificationservice/application/query"
	"github.com/janicaleksander/cloud/policyverificationservice/infrastructure/messaging"
	dynamoDB "github.com/janicaleksander/cloud/policyverificationservice/infrastructure/tableDB"
	"github.com/janicaleksander/cloud/policyverificationservice/persistance"
	"github.com/janicaleksander/cloud/policyverificationservice/presentation"
	"github.com/janicaleksander/cloud/policyverificationservice/presentation/router"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	db, err := dynamoDB.NewTableDB()
	if err != nil {
		panic(err)
	}
	err = db.Migrate()
	if err != nil {
		panic(err)
	}
	rabbit, err := rabbitmq.NewRabbitMQ()
	if err != nil {
		panic(err)
	}
	publisher := rabbitmq.NewPublisher(rabbit)

	policyRepository := persistance.NewPolicyRepository(db)

	createPolicyHandler := command.NewCreatePolicyCommandHandler(policyRepository)
	deletePolicyHandler := command.NewDeletePolicyCommandHandler(policyRepository)
	updatePolicyHandler := command.NewUpdatePolicyCommandHandler(policyRepository)

	err = createPolicyHandler.SelfRegister()
	if err != nil {
		slog.Error("Error registering create policy command handler", "error", err)
		panic(err)
	}
	err = deletePolicyHandler.SelfRegister()
	if err != nil {
		slog.Error("Error registering delete policy command handler", "error", err)
		panic(err)
	}
	err = updatePolicyHandler.SelfRegister()
	if err != nil {
		slog.Error("Error registering update policy command handler", "error", err)
		panic(err)
	}

	getPolicyHandler := query.NewGetPolicyQueryHandler(policyRepository)
	getPoliciesHandler := query.NewGetPoliciesQueryHandler(policyRepository)
	checkPolicyHandler := query.NewCheckPolicyQueryHandler(policyRepository, publisher)

	err = getPolicyHandler.SelfRegister()
	if err != nil {
		slog.Error("Error registering get policy query handler", "error", err)
		panic(err)
	}
	err = getPoliciesHandler.SelfRegister()
	if err != nil {
		slog.Error("Error registering get policies query handler", "error", err)
		panic(err)
	}
	err = checkPolicyHandler.SelfRegister()
	if err != nil {
		slog.Error("Error registering check policy query handler", "error", err)
		panic(err)
	}

	policyController := presentation.NewPolicyController()

	policyEventHandler := messaging.NewPolicyEventHandler()
	go policyEventHandler.Run(rabbit)

	r := router.NewRouter(policyController)
	log.Println("serving on 8081")
	err = http.ListenAndServe("localhost:8081", r)
	if err != nil {
		slog.Error("Error running http: ", err)
		panic(err)
	}
}
