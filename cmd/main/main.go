package main

import (
	_ "github.com/a-h/templ"
	"github.com/joho/godotenv"

	"space-traders/service/api"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	err = api.NewAPI().Start()
	if err != nil {
		panic(err)
	}
}
