package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {

	data := envelope{
		"status": "available",
		"system_info": map[string]string{
			"enviroment": app.config.env,
			"version":    version,
		},
	}

	err := app.writeJSON(w, 200, data, r.Header)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "server error", http.StatusInternalServerError)
	}

}
