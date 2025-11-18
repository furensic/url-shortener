package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type Middleware func(next http.Handler) http.Handler

func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		next.ServeHTTP(w, r)

		elapsedTime := time.Since(startTime)
		log.Printf("[%s] Path: %s Elapsed: %s\n", r.Method, r.URL.Path, elapsedTime)
	})
}

func (app *application) serv() error {
	log.Print("Creating http server")
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", app.port),
		Handler: Log(app.routes()),
	}

	log.Print("Starting http server on ", server.Addr)

	return server.ListenAndServe()
}
