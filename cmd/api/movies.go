package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ggarber42/greenlight/internal/data"
	"github.com/ggarber42/greenlight/internal/validator"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
	} else {

		movie := &data.Movie{
			Title:   input.Title,
			Year:    input.Year,
			Runtime: input.Runtime,
			Genres:  input.Genres,
		}
		v := validator.New()

		if data.ValidateMovie(v, movie); !v.Valid() {
			app.failedValidationResponse(w, r, v.Errors)
		} else {
			err := app.models.Movies.Insert(movie)
			if err != nil {
				app.serverErrorResponse(w, r, err)
			} else {
				headers := make(http.Header)
				headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))
				err := app.writeJSON(w, http.StatusCreated, envelope{"movie": movie}, headers)
				if err != nil {
					app.serverErrorResponse(w, r, err)
				}
			}
		}

	}

}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := readIDParam(r)

	if err != nil {
		app.notFoundResponse(w, r)
	} else {

		movie, err := app.models.Movies.Get(id)

		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.notFoundResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
		} else {
			err := app.writeJSON(w, 200, envelope{"movie": movie}, r.Header)

			if err != nil {
				app.serverErrorResponse(w, r, err)
			}
		}

	}
}

func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := readIDParam(r)

	if err != nil {
		app.notFoundResponse(w, r)
	} else {

		movie, err := app.models.Movies.Get(id)

		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.notFoundResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
		} else {

			var input struct {
				Title   string       `json:"title"`
				Year    int32        `json:"year"`
				Runtime data.Runtime `json:"runtime"`
				Genres  []string     `json:"genres"`
			}

			err := app.readJSON(w, r, &input)
			if err != nil {
				app.errorResponse(w, r, http.StatusBadRequest, err.Error())
			} else {

				movie.Title = input.Title
				movie.Year = input.Year
				movie.Runtime = input.Runtime
				movie.Genres = input.Genres

				v := validator.New()

				if data.ValidateMovie(v, movie); !v.Valid() {
					app.failedValidationResponse(w, r, v.Errors)
				} else {
					err := app.models.Movies.Update(movie, id)

					if err != nil {
						app.serverErrorResponse(w, r, err)
					} else {
						err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
						if err != nil {
							app.serverErrorResponse(w, r, err)
						}
					}

				}
			}
		}

	}
}
