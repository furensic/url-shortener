package main

import (
	"log/slog"
	"net/http"

	"codeberg.org/Kassiopeia/url-shortener/internal/models"
)

func (app *application) basicAuthMiddleware(next http.Handler) http.Handler {
	app.logger.Debug("basic auth middleware")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		app.logger.Debug("basic auth middleware")
		if !ok {
			slog.Info("No authorization header set")
			http.Error(w, "No authorization header", http.StatusUnauthorized)
			return
		}

		loginPayload := models.LoginUserPayload{
			Username: username,
			Password: password,
		}

		ok, err := app.service.UserService.VerifyCredentials(loginPayload)
		if err != nil {
			slog.Error(err.Error())
		}

		if ok {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

	})
}

func (app *application) mountRoutes() http.Handler {
	app.logger.Debug("Creating public new mux")
	publicMux := http.NewServeMux()
	// health
	app.logger.Debug("Mounting GET /health")
	publicMux.Handle("GET /health", app.basicAuthMiddleware(http.HandlerFunc(app.GetHealthHandler))) // ?

	// shortened_uri
	app.logger.Debug("Mounting POST /")
	publicMux.HandleFunc("POST /", app.CreateShortenedUri) // ?
	app.logger.Debug("Mounting GET /{id}")
	publicMux.HandleFunc("GET /{id}", app.GetShortUriById)
	publicMux.HandleFunc("GET /redirect/{id}", app.GetShortUriByIdRedirect) // ?

	// Auth
	publicMux.HandleFunc("POST /auth/register", app.RegisterUser)
	publicMux.HandleFunc("POST /auth/login", app.LoginUser)
	publicMux.HandleFunc("PUT /auth/password", app.UpdatePassword)

	// Users
	publicMux.HandleFunc("GET /user/{username}", app.GetUserByName)
	publicMux.HandleFunc("PUT /user/{id}", app.UpdateUserExtension)

	// root router
	app.logger.Debug("Creating root mux")
	rootMux := http.NewServeMux()

	app.logger.Debug("Add /v1/ handle to root mux")
	rootMux.Handle("/v1/", http.StripPrefix("/v1", logMiddleware(publicMux)))

	return rootMux
}
