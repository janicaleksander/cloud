package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/janicaleksander/cloud/common/rabbitmq"
	"github.com/janicaleksander/cloud/valuationservice/application/command"
	"github.com/janicaleksander/cloud/valuationservice/application/query"
	"github.com/janicaleksander/cloud/valuationservice/infrastructure"
	"github.com/janicaleksander/cloud/valuationservice/infrastructure/ai"
	"github.com/janicaleksander/cloud/valuationservice/infrastructure/messaging"
	"github.com/janicaleksander/cloud/valuationservice/persistence"
	"github.com/janicaleksander/cloud/valuationservice/presentation"
	"github.com/janicaleksander/cloud/valuationservice/presentation/router"
)

func main() {
	/*	err := godotenv.Load()
		if err != nil {
			panic(err)
		}*/
	db, err := infrastructure.NewDB()
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&persistence.ValuationModel{}, &persistence.PartModel{})
	if err != nil {
		panic(err)
	}

	rabbit, err := rabbitmq.NewRabbitMQ()
	publisher := rabbitmq.NewPublisher(rabbit)
	newMockDamageDetectro := ai.NewMockDamageDetector()
	if err != nil {
		panic(err)
	}
	valuationRepository := persistence.NewValuationRepository(db)

	calculateValuationCommand := command.NewCalculateValuationCommandHandler(valuationRepository, newMockDamageDetectro, publisher)
	createValuationCommand := command.NewCreateValuationCommandHandler(valuationRepository)
	deleteValuationCommand := command.NewDeleteValuationCommandHandler(valuationRepository)
	updateValuationCommand := command.NewUpdateValuationCommandHandler(valuationRepository)

	_ = calculateValuationCommand.SelfRegister()
	_ = createValuationCommand.SelfRegister()
	_ = deleteValuationCommand.SelfRegister()
	_ = updateValuationCommand.SelfRegister()

	getValuationQuery := query.NewGetValuationQueryHandler(valuationRepository)
	getValuationsQuery := query.NewGetValuationsQueryHandler(valuationRepository)

	_ = getValuationQuery.SelfRegister()
	_ = getValuationsQuery.SelfRegister()

	valuationController := presentation.NewValuationController()
	valuationHandler := messaging.NewValuationEventHandler()
	valuationHandler.Run(rabbit)
	r := router.NewRouter(valuationController)
	slog.Info("Starting valuation service on port ", os.Getenv("APP_PORT"))
	err = http.ListenAndServe("localhost:8082", r)
	if err != nil {
		slog.Error("Error running http: ", err)
		panic(err)
	}

}

//TODO sprawdzic czy importy sa z dobrych paczke i zamiana teogb api na presentation

//Idempotentność
//
//Co się stanie jeśli ten sam event zostanie przetworzony dwa razy? Np. przy Nack + requeue handler dostanie wiadomość ponownie.
//Czy CheckUserPolicy, ChangeClaimStatus itp. są bezpieczne przy wielokrotnym wywołaniu z tymi samymi danymi?
