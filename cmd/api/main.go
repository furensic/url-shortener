package main

import (
	"log"
)

type application struct {
	port int
}

func main() {
	log.Print("Starting url shortener service")

	app := &application{
		port: 8080,
	}

	if err := app.serv(); err != nil {
		log.Fatal("Error listening and serving:", err.Error())
	}
}
