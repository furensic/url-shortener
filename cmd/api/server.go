package main

import (
	"fmt"
	"log"
	"net/http"
)

func (app *application) serv() error {
	log.Print("Creating http server")
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", app.port),
		Handler: app.routes(),
	}

	log.Print("Starting http server on ", server.Addr)

	return server.ListenAndServe()
}
