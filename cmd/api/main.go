package main

import (
	"log"

	"codeberg.org/Kassiopeia/url-shortener/internal/database"
)

type application struct {
	port   int
	models database.Models
}

func main() {
	log.Print("Starting url shortener service")

	app := &application{
		port:   8080,
		models: database.NewModels(),
	}

	if err := app.serv(); err != nil {
		log.Fatal("Error listening and serving:", err.Error())
	}
}
