package main

import (
	"log"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	log.Print("Mounting handlers")
	mux.HandleFunc("GET /health", app.getHealth)

	return mux
}
