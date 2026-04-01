package main

import (
	"fmt"

	"github.com/janicaleksander/cloud/common/rabbitmq"
	"github.com/janicaleksander/cloud/notificationservice/infrastructure"
	"github.com/janicaleksander/cloud/notificationservice/infrastructure/messaging"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	_, err = infrastructure.NewDB()
	if err != nil {
		panic(err)
	}
	notificationEventHandler := messaging.NewNotificationHandler()
	rabbit, err := rabbitmq.NewRabbitMQ()
	if err != nil {
		panic(err)
	}
	err = notificationEventHandler.Run(rabbit)
	if err != nil {
		panic(err)
	}
	fmt.Println("Notification service is running...")
	select {}
}
