package main

import (
	"net/http"

	"codeberg.org/Kassiopeia/url-shortener/cmd/api/handlers"
)

func (app *application) mountRoutes(h *handlers.Handler) http.Handler {
	app.logger.Debug("Creating public new mux")
	publicMux := http.NewServeMux()

	// health
	app.logger.Debug("Mounting GET /health")
	publicMux.HandleFunc("GET /health", h.GetHealthHandler) // ?

	// shortened_uri
	app.logger.Debug("Mounting POST /")
	publicMux.HandleFunc("POST /", h.CreateShortenedUri) // ?
	app.logger.Debug("Mounting GET /{id}")
	publicMux.HandleFunc("GET /{id}", h.GetShortenedUriById) // ?

	publicMux.HandleFunc("POST /auth/register", h.RegisterUser)

	// root router
	app.logger.Debug("Creating root mux")
	rootMux := http.NewServeMux()

	app.logger.Debug("Add /v1/ handle to root mux")
	rootMux.Handle("/v1/", http.StripPrefix("/v1", logMiddleware(publicMux)))

	return rootMux
}
