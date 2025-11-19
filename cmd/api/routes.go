package main

import (
	"log"
	"net/http"

	"codeberg.org/Kassiopeia/url-shortener/cmd/api/handlers"
)

func (app *application) mountRoutes(h *handlers.Handler) http.Handler {
	log.Print("Mounting handlers")

	publicMux := http.NewServeMux()
	// health
	publicMux.HandleFunc("GET /health", h.GetHealthHandler) // ?
	// shortened_uri
	publicMux.HandleFunc("POST /", h.CreateShortenedUri)     // ?
	publicMux.HandleFunc("GET /{id}", h.GetShortenedUriById) // ?

	// root router
	rootMux := http.NewServeMux()
	rootMux.Handle("/v1/", http.StripPrefix("/v1", logMiddleware(publicMux)))

	return rootMux
}
