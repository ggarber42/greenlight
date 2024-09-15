package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ggarber42/greenlight/internal/data"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string   `json:"title"`
		Year    int32    `json:"year"`
		Runtime data.Runtime    `json:"runtime"`
		Genres  []string `json:"genres"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
	} else {
		fmt.Fprintf(w, "%+v\n", input)
	}

}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := readIDParam(r)

	if err != nil {
		app.notFoundResponse(w, r)
	} else {

		movie := data.Movie{
			ID:        id,
			CreatedAt: time.Now(),
			Title:     "Casablanca",
			Runtime:   102,
			Genres:    []string{"drama", "romance", "war"},
			Version:   1,
		}

		err = app.writeJSON(w, 200, envelope{"movie": movie}, r.Header)

		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
	}
}
