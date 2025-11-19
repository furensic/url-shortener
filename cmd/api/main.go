package main

import (
	"context"
	"log"

	"codeberg.org/Kassiopeia/url-shortener/internal/database"

	"github.com/jackc/pgx/v5"
)

type application struct {
	port   int
	models database.Models
}

func main() {
	log.Print("Starting url shortener service")

	db, err := pgx.Connect(context.Background(), "postgres://svc:password@postgres:5432/url_shortener")
	if err != nil {
		log.Fatal("Error opening database connection: ", err.Error())
	}

	app := &application{
		port:   8080,
		models: database.NewModels(db),
	}

	if err := app.serv(); err != nil {
		log.Fatal("Error listening and serving:", err.Error())
	}
}
