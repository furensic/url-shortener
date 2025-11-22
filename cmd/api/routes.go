package main

import (
	"log/slog"
	"net/http"
	"time"

	"codeberg.org/Kassiopeia/url-shortener/internal/models"
)

func (app *application) basicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
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

type Middleware func(next http.Handler) http.Handler

func (app *application) logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		next.ServeHTTP(w, r)

		elapsedTime := time.Since(startTime)
		app.logger.WithGroup("Request").Info("Logging request", slog.String("method", r.Method), slog.String("path", r.URL.Path), slog.Duration("request time", elapsedTime))
		// not sure how to replace this standard logger with my application logger?
	})
}

func (app *application) mountRoutes() http.Handler {
	publicMux := http.NewServeMux()
	// health
	publicMux.Handle("GET /health", app.basicAuthMiddleware(http.HandlerFunc(app.GetHealthHandler))) // ?

	// shortened_uri
	publicMux.HandleFunc("POST /", app.CreateShortenedUri) // ?
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
	rootMux := http.NewServeMux()

	rootMux.Handle("/v1/", http.StripPrefix("/v1", app.logMiddleware(publicMux)))

	return rootMux
}
