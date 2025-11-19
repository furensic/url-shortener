package main

import (
	"log"
	"net/http"
)

func (app *application) mountRoutes() http.Handler {
	log.Print("Mounting handlers")

	handlers := &Handler{} // ?

	publicMux := http.NewServeMux()
	// health
	publicMux.HandleFunc("GET /health") // ?
	// shortened_uri
	publicMux.HandleFunc("POST /", app.createShortenedUri)     // ?
	publicMux.HandleFunc("GET /{id}", app.getShortenedUriById) // ?

	// root router
	rootMux := http.NewServeMux()
	rootMux.Handle("/v1/", http.StripPrefix("/v1", logMiddleware(publicMux)))

	return rootMux
}
