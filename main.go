package main

import (
	"os"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		return
	}
	conn, err := amqp.Dial(os.Getenv("AMQP_URL"))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
}
