package main

import (
	"fmt"
	"net/http"
)

func (app *application) GetHealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("Application running!\nURL Query: %v", r.URL.Query())))

}
