package main

import (
	"log"

	_ "github.com/a-h/templ"
	"github.com/joho/godotenv"

	"space-traders/service/api"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Panicf("Error loading .env file: %v", err)
	}

	err = api.NewAPI().Start()
	if err != nil {
		log.Panicf("Error starting API: %v", err)
	}
}
