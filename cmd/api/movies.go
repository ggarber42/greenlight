package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ggarber42/greenlight/internal/data"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new movie")
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := readIDParam(r)

	if err != nil {
		http.NotFound(w, r)
	} else {

		movie := data.Movie{
			ID:        id,
			CreatedAt: time.Now(),
			Title:     "Casablanca",
			Runtime:   102,
			Genres:    []string{"drama", "romance", "war"},
			Version:   1,
		}

		err = app.writeJSON(w, 200, movie, r.Header)

		if err != nil {
			app.logger.Println(err)
			http.Error(w, "serve encontured an error", http.StatusInternalServerError)
		}
	}
}
