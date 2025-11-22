package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) GetUserByName(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	if username == "" {
		http.Error(w, "no username found in path arguments", http.StatusBadRequest)
		return
	}

	user, err := app.service.UserService.GetByUsername(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")
	enc := json.NewEncoder(w)
	enc.Encode(user)
}

func (app *application) UpdateUserExtension(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Unimplemented endpoint", http.StatusNotImplemented)
}
