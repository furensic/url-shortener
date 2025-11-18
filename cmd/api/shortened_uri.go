package main

import (
	"log"
	"net/http"

	"codeberg.org/Kassiopeia/url-shortener/internal/database"
)

func (app *application) createShortenedUri(w http.ResponseWriter, r *http.Request) {
	log.Print("Received POST / request. Handler function createShortenedUri")
	if err := app.models.ShortenedUri.Create(&database.ShortenedUri{}); err != nil {
		log.Fatal("Error creating shortenedUri: ", err.Error())
	}
}

func (app *application) getShortenedUriById(w http.ResponseWriter, r *http.Request) {
	log.Print("Received GET / request. Handler function getShortenedUriById")
	w.Write([]byte("Not implemented!\n"))
}
