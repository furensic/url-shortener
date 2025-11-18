package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

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
	param_id := r.PathValue("id")
	id, err := strconv.ParseInt(param_id, 10, 64)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Error parsing path parameter: %v\n", err.Error())))
		return
	}
	w.Write([]byte(fmt.Sprintf("id: %d\n", id)))
}
