package main

import (
	"encoding/json"
	"net/http"

	"codeberg.org/Kassiopeia/url-shortener/internal/models"
)

func (app *application) RegisterUser(w http.ResponseWriter, r *http.Request) {
	newRequest := models.RegisterUserPayload{}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&newRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := app.service.UserService.Create(newRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")
	enc := json.NewEncoder(w)
	enc.Encode(user)
}

func (app *application) LoginUser(w http.ResponseWriter, r *http.Request) {
	newRequest := models.LoginUserPayload{}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&newRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := app.service.UserService.VerifyCredentials(newRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")
	enc := json.NewEncoder(w)
	enc.Encode(user)
}
