package main

import (
	"log"
	"net/http"
)

func main() {
	log.Print("Starting url shortener service")
	mux := http.NewServeMux()

	log.Print("Mounting handlers")
	mux.HandleFunc("GET /health", getHealth)

	log.Print("Creating http server")
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Print("Starting http server on ", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Error listening and serving:", err.Error())
	}
}

func getHealth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Application running!\n"))
}
