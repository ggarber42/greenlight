package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {

	data := map[string]string{
		"status":     "available",
		"enviroment": app.config.env,
		"version": version,
	}

	json, err := json.Marshal(data)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "server error", http.StatusInternalServerError)
	}  else {

		w.Header().Set("Content-Type", "application/json")
	
		w.Write([]byte(json))
	}

}
