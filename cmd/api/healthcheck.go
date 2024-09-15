package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {

	data := map[string]string{
		"status":     "available",
		"enviroment": app.config.env,
		"version": version,
	}

	err := app.writeJSON(w, 200, data, r.Header)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "server error", http.StatusInternalServerError)
	} 

}
