package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"codeberg.org/Kassiopeia/url-shortener/internal/models"
	"codeberg.org/Kassiopeia/url-shortener/internal/service"
	"github.com/jackc/pgx/v5"
)

type Handler struct {
	UriService service.ShortenedUriService
}

type CreateShortenedUriRequest struct {
	OriginUri string `json:"origin_uri"`
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) getHealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Application running!\n"))
}

func (h *Handler) createShortenedUri(w http.ResponseWriter, r *http.Request) {
	shortenedUri := models.ShortenedUri{}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&shortenedUri); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tmp, err := h.UriService.Create(&shortenedUri)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")
	enc := json.NewEncoder(w)
	enc.Encode(tmp)
}

func (h *Handler) getShortenedUriById(w http.ResponseWriter, r *http.Request) {
	param_id := r.PathValue("id")
	id, err := strconv.Atoi(param_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	shortenedUri, err := h.UriService.GetById(id)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Print("Redirecting request to:", shortenedUri.OriginUri)
	http.Redirect(w, r, shortenedUri.OriginUri, http.StatusPermanentRedirect)
}
