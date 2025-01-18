package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"greenlight.Arkitecth.github/internal/data"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new movie")
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	movie := data.Movie{
		ID: id,
		CreatedAt: time.Now(),
		Title: "Casablanca",
		Runtime: 102,
		Genres: []string{"drama", "romance", "war"},
		Version: 1,
	}
	err  = app.writeJSON(w, http.StatusOK, movie, nil) 
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}

func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}
	return id, nil

}