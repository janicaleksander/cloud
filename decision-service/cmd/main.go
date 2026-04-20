package main

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/janicaleksander/cloud/common/rabbitmq"
	"github.com/janicaleksander/cloud/decisionservice/application/command"
	"github.com/janicaleksander/cloud/decisionservice/application/query"
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

	deleteDecisionCommand := command.NewDeleteDecisionCommandHandler(decisionRepository)
	prepareDecisionCommand := command.NewPrepareDecisionCommandHandler(decisionRepository)
	updateDecisionStateCommand := command.NewUpdateDecisionStateCommandHandler(decisionRepository, publisher)
	updateEmpCommand := command.NewUpdateEmpCommandHandler(decisionRepository)

	_ = deleteDecisionCommand.SelfRegister()
	_ = prepareDecisionCommand.SelfRegister()
	_ = updateDecisionStateCommand.SelfRegister()
	_ = updateEmpCommand.SelfRegister()

	getDecisionQuery := query.NewGetDecisionQueryHandler(decisionRepository)
	getDecisionsQuery := query.NewGetDecisionsQueryHandler(decisionRepository)
	getWaitingDecisions := query.NewGetWaitingDecisionsQueryHandler(decisionRepository)
	_ = getDecisionQuery.SelfRegister()
	_ = getDecisionsQuery.SelfRegister()
	_ = getWaitingDecisions.SelfRegister()

	decisionController := presentation.NewDecisionController()
	decisionEventHandler := messaging.NewDecisionEventHandler()
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
