package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type config struct {
	port              int
	writeTimeout      time.Duration
	readTimeout       time.Duration
	readHeaderTimeout time.Duration
	idleTimeout       time.Duration

	tls            bool
	certificate    string
	certificateKey string
}

type Middleware func(next http.Handler) http.Handler

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		next.ServeHTTP(w, r)

		elapsedTime := time.Since(startTime)
		log.Printf("[%s] Path: %s Elapsed: %s\n", r.Method, r.URL.Path, elapsedTime)
		// not sure how to replace this standard logger with my application logger?
	})
}

func (app *application) serveHTTP() error {
	app.logger.Debug("Mounting routes from handler")
	router := app.mountRoutes()

	server := http.Server{
		Addr:              fmt.Sprintf(":%d", app.config.port),
		Handler:           router,
		WriteTimeout:      app.config.writeTimeout,
		ReadTimeout:       app.config.readTimeout,
		ReadHeaderTimeout: app.config.readHeaderTimeout,
		IdleTimeout:       app.config.idleTimeout,
	}
	app.logger.Debug("Created new http.Server config " + fmt.Sprintf("%+v", &server))

	app.logger.Debug("Starting http server on " + server.Addr)

	if app.config.tls && app.config.certificate != "" && app.config.certificateKey != "" {
		return server.ListenAndServeTLS(app.config.certificate, app.config.certificateKey)
	} else {
		return server.ListenAndServe()
	}
}
