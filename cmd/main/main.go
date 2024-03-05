package main

import (
	_ "github.com/a-h/templ"
	"github.com/joho/godotenv"

	"space-traders/service/api"
)

func main() {
	godotenv.Load()

	err := api.NewAPI().Start()
	if err != nil {
		panic(err)
	}
}
