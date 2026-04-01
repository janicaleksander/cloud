package main

import (
	"log"
	"log/slog"

	"github.com/janicaleksander/cloud/common/rabbitmq"
	"github.com/janicaleksander/cloud/notificationservice/application"
	"github.com/janicaleksander/cloud/notificationservice/infrastructure"
	"github.com/janicaleksander/cloud/notificationservice/infrastructure/messaging"
	"github.com/janicaleksander/cloud/notificationservice/persistance"
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

	err = db.AutoMigrate(&persistance.NotificationReceiverModel{}, &persistance.NotificationModel{})
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
	notificationRepository := persistance.NewNotificationRepository(db)
	emailService := application.NewEmailService("notifications@insurance.com")
	_ = application.NewNotificationService(notificationRepository)

	notificationEventHandler := messaging.NewNotificationHandler(emailService)
	err = notificationEventHandler.Run(rabbit)
	if err != nil {
		slog.Error("Error running notification event handler", "error", err)
		panic(err)
	}

	log.Println("Notification service is running...")
	select {}
}
