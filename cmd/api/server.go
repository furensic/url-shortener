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
}

type Middleware func(next http.Handler) http.Handler

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		next.ServeHTTP(w, r)

		elapsedTime := time.Since(startTime)
		log.Printf("[%s] Path: %s Elapsed: %s\n", r.Method, r.URL.Path, elapsedTime)
	})
}

func (app *application) serveHTTP() error {
	log.Print("Creating http server")
	server := http.Server{
		Addr:              fmt.Sprintf(":%d", app.config.port),
		Handler:           app.mountRoutes(),
		WriteTimeout:      app.config.writeTimeout,
		ReadTimeout:       app.config.readTimeout,
		ReadHeaderTimeout: app.config.readHeaderTimeout,
		IdleTimeout:       app.config.idleTimeout,
	}

	log.Print("Starting http server on ", server.Addr)

	return server.ListenAndServe()
}
