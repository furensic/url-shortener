package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"codeberg.org/Kassiopeia/url-shortener/internal/models"
	"codeberg.org/Kassiopeia/url-shortener/internal/repository"
	"codeberg.org/Kassiopeia/url-shortener/internal/service"
)

type CreateShortenedUriRequest struct {
	OriginUri string `json:"origin_uri"`
}

type Handler struct {
	ShortenerService *service.ShortenerService
	UserService      *service.UserService
}

func (h *Handler) GetHealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Application running!\n"))
}

func (h *Handler) CreateShortenedUri(w http.ResponseWriter, r *http.Request) {
	shortenedUri := models.ShortenedUri{}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&shortenedUri); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tmp, err := h.ShortenerService.Create(shortenedUri)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")
	enc := json.NewEncoder(w)
	enc.Encode(tmp)
}

func (h *Handler) GetShortenedUriById(w http.ResponseWriter, r *http.Request) {
	param_id := r.PathValue("id")
	id, err := strconv.Atoi(param_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Print("before h.ShortenerService.GetById(id)")

	shortenedUri, err := h.ShortenerService.GetById(id)
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

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	newRequest := models.RegisterUserPayload{}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&newRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.UserService.Create(newRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")
	enc := json.NewEncoder(w)
	enc.Encode(user)
}

func (h *Handler) GetUserByName(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	if username == "" {
		http.Error(w, "no username found in path arguments", http.StatusBadRequest)
		return
	}

	user, err := h.UserService.GetByUsername(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")
	enc := json.NewEncoder(w)
	enc.Encode(user)
}
