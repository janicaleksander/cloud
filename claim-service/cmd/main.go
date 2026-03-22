package main

import (
	"github.com/janicaleksander/cloud/claimservice/infrastructure"
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
}
