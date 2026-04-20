package main

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/janicaleksander/cloud/common/rabbitmq"
	"github.com/janicaleksander/cloud/notificationservice/application"
	"github.com/janicaleksander/cloud/notificationservice/application/command"
	"github.com/janicaleksander/cloud/notificationservice/application/query"
	"github.com/janicaleksander/cloud/notificationservice/infrastructure"
	"github.com/janicaleksander/cloud/notificationservice/infrastructure/messaging"
	"github.com/janicaleksander/cloud/notificationservice/persistence"
	"github.com/janicaleksander/cloud/notificationservice/presentation"
	"github.com/janicaleksander/cloud/notificationservice/presentation/router"
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

	err = db.AutoMigrate(&persistence.NotificationReceiverModel{}, &persistence.NotificationModel{})
	if err != nil {
		slog.Error("Error migrating database", "error", err)
		panic(err)
	}

	rabbit, err := rabbitmq.NewRabbitMQ()
	if err != nil {
		slog.Error("Error connecting to RabbitMQ", "error", err)
		panic(err)
	}

	// Initialize repositories and services
	notificationRepository := persistence.NewNotificationRepository(db)
	emailService := application.NewEmailService("notifications@insurance.com")

	createNotificationCommand := command.NewCreateNotificationCommandHandler(notificationRepository)
	createNotificationReceiver := command.NewCreateNotificationReceiverCommandHandler(notificationRepository)
	deleteNotificationCommand := command.NewDeleteNotificationCommandHandler(notificationRepository)

	_ = createNotificationCommand.SelfRegister()
	_ = createNotificationReceiver.SelfRegister()
	_ = deleteNotificationCommand.SelfRegister()

	getEmailByClaimIDQuery := query.NewGetEmailByClaimIDQueryHandler(notificationRepository)
	getNotificationQuery := query.NewGetNotificationQueryHandler(notificationRepository)
	getNotificationsQuery := query.NewGetNotificationsQueryHandler(notificationRepository)
	getNotificationsForClaimIDQuery := query.NewGetNotificationsForClaimIDQueryHandler(notificationRepository)

	_ = getEmailByClaimIDQuery.SelfRegister()
	_ = getNotificationQuery.SelfRegister()
	_ = getNotificationsQuery.SelfRegister()
	_ = getNotificationsForClaimIDQuery.SelfRegister()

	notificationController := presentation.NewNotificationController()
	notificationEventHandler := messaging.NewNotificationHandler(emailService)
	err = notificationEventHandler.Run(rabbit)
	if err != nil {
		slog.Error("Error running notification event handler", "error", err)
		panic(err)
	}
	r := router.NewRouter(notificationController)
	log.Println("Notification service is running...")

	err = http.ListenAndServe(":8085", r)

	if err != nil {
		slog.Error("Error starting HTTP server", "error", err)
		panic(err)
	}

}
