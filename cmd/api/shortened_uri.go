package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"codeberg.org/Kassiopeia/url-shortener/internal/database"
	"github.com/jackc/pgx/v5"
)

func (app *application) createShortenedUri(w http.ResponseWriter, r *http.Request) {
	log.Print("Received POST / request. Handler function createShortenedUri")

	shortenedUri := database.ShortenedUri{}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&shortenedUri); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := app.models.ShortenedUri.Create(&shortenedUri); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (app *application) getShortenedUriById(w http.ResponseWriter, r *http.Request) {
	log.Print("Received GET / request. Handler function getShortenedUriById")
	param_id := r.PathValue("id")
	id, err := strconv.Atoi(param_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	shortenedUri, err := app.models.ShortenedUri.GetById(id)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	enc := json.NewEncoder(w)
	enc.Encode(shortenedUri)
}
