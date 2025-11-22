package main

import (
	"fmt"
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

func (app *application) serveHTTP() error {
	router := app.mountRoutes()

	server := http.Server{
		Addr:              fmt.Sprintf(":%d", app.config.port),
		Handler:           router,
		WriteTimeout:      app.config.writeTimeout,
		ReadTimeout:       app.config.readTimeout,
		ReadHeaderTimeout: app.config.readHeaderTimeout,
		IdleTimeout:       app.config.idleTimeout,
	}

	if app.config.tls && app.config.certificate != "" && app.config.certificateKey != "" {
		return server.ListenAndServeTLS(app.config.certificate, app.config.certificateKey)
	} else {
		return server.ListenAndServe()
	}
}
