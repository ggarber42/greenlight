package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

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

			if r.Header.Get("X-Expected-Version") != "" {
				if strconv.FormatInt(int64(movie.Version), 32) != r.Header.Get("X-Expected-Version") {
					app.editConflictResponse(w, r)
				}
			} else {

			}

			var input struct {
				Title   *string       `json:"title"`
				Year    *int32        `json:"year"`
				Runtime *data.Runtime `json:"runtime"`
				Genres  []string      `json:"genres"`
			}

			err := app.readJSON(w, r, &input)
			if err != nil {
				app.errorResponse(w, r, http.StatusBadRequest, err.Error())
			} else {

				if input.Title != nil {
					movie.Title = *input.Title
				}

				if input.Year != nil {
					movie.Year = *input.Year
				}

				if input.Runtime != nil {
					movie.Runtime = *input.Runtime
				}

				if input.Genres != nil {
					movie.Genres = input.Genres
				}

				v := validator.New()

				if data.ValidateMovie(v, movie); !v.Valid() {
					app.failedValidationResponse(w, r, v.Errors)
				} else {
					err := app.models.Movies.Update(movie)

					if err != nil {
						switch {
						case errors.Is(err, data.ErrEditConflict):
							app.editConflictResponse(w, r)
						default:
							app.serverErrorResponse(w, r, err)
						}
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

func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
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
			err := app.models.Movies.Delete(movie.ID)

			if err != nil {
				app.serverErrorResponse(w, r, err)
			} else {
				err = app.writeJSON(w, http.StatusOK, envelope{"message": "movie successfully deleted"}, nil)
				if err != nil {
					app.serverErrorResponse(w, r, err)
				}
			}
		}
	}
}

func (app *application) listMoviesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string
		Genres []string
		data.Filters
	}

	v := validator.New()
	qs := r.URL.Query()
	input.Title = app.readString(qs, "title", "")
	input.Genres = app.readCSV(qs, "genres", []string{})

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
	} else {

		movies, metadata, err := app.models.Movies.GetAll(input.Title, input.Genres, input.Filters)

		if err != nil {
			app.serverErrorResponse(w, r, err)
		} else {
			err := app.writeJSON(w, http.StatusOK, envelope{"metadata": metadata, "movies": movies}, nil)

			if err != nil {
				app.serverErrorResponse(w, r, err)
			}
		}

	}

}
