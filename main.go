package main

import (
	"fmt"
	"os"
	"reflect"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

type eventevent struct {
}

func main() {
	t := reflect.TypeOf(eventevent{})
	fmt.Println(t.Name())
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
