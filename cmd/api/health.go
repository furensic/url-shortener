package main

import "net/http"

func (app *application) getHealth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Application running!\n"))
}
