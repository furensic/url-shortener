package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"codeberg.org/Kassiopeia/url-shortener/internal/models"
	"codeberg.org/Kassiopeia/url-shortener/internal/repository"
)

type CreateShortenedUriRequest struct {
	OriginUri string `json:"origin_uri"`
}

func (app *application) CreateShortenedUri(w http.ResponseWriter, r *http.Request) {
	shortenedUri := models.ShortenedUri{}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&shortenedUri); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tmp, err := app.service.ShortenerService.Create(shortenedUri) // app.service.ShortenerService undefined (type Services has no field or method ShortenerService)go-staticcheck
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")
	enc := json.NewEncoder(w)
	enc.Encode(tmp)
}

func (app *application) GetShortUriByIdRedirect(w http.ResponseWriter, r *http.Request) {
	param_id := r.PathValue("id")
	id, err := strconv.Atoi(param_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Print("before h.ShortenerService.GetById(id)")

	shortenedUri, err := app.service.ShortenerService.GetById(id) // app.service.ShortenerService undefined (type Services has no field or method ShortenerService) (compile)go-staticcheck
	if err != nil {
		if err == repository.ErrShortenedUriNotFound {
			log.Print("ErrNoRows GetById")
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Print("Redirecting request to:", shortenedUri.OriginUri)
	http.Redirect(w, r, shortenedUri.OriginUri, http.StatusPermanentRedirect)
}

func (app *application) GetShortUriById(w http.ResponseWriter, r *http.Request) {
	param_id := r.PathValue("id")
	id, err := strconv.Atoi(param_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Print("before h.ShortenerService.GetById(id)")

	shortenedUri, err := app.service.ShortenerService.GetById(id) // app.service.ShortenerService undefined (type Services has no field or method ShortenerService) (compile)go-staticcheck
	if err != nil {
		if err == repository.ErrShortenedUriNotFound {
			log.Print("ErrNoRows GetById")
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
