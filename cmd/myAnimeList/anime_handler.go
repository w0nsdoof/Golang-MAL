package main

import (
	"errors"
	"log"
	"net/http"

	"final-project/pkg/model"
	"final-project/pkg/validator"
)

func (app *application) createAnimeHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Rating float64 `json:"rating"`
		Title  string  `json:"title"`
		Genres string  `json:"genres"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		log.Println(err)
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}

	anime := &model.Anime{
		Rating: input.Rating,
		Title:  input.Title,
		Genres: input.Genres,
	}

	err = app.models.Animes.Insert(anime)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"anime": anime}, nil)
}

func (app *application) getAnimesList(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string
		model.Filters
	}
	v := validator.New()
	qs := r.URL.Query()

	input.Title = app.readStrings(qs, "title", "")

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readStrings(qs, "sort", "id")

	input.Filters.SortSafeList = []string{
		"id", "title",
		"-id", "-title",
	}

	if model.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	animes, metadata, err := app.models.Animes.GetAll(input.Title, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"animes": animes, "metadata": metadata}, nil)
}

func (app *application) getAnimeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	anime, err := app.models.Animes.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"anime": anime}, nil)
}

func (app *application) updateAnimeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	anime, err := app.models.Animes.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Rating *float64 `json:"rating"`
		Title  *string  `json:"title"`
		Genres *string  `json:"genres"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Rating != nil {
		anime.Rating = *input.Rating
	}
	if input.Title != nil {
		anime.Title = *input.Title
	}
	if input.Genres != nil {
		anime.Genres = *input.Genres
	}

	v := validator.New()

	if model.ValidateAnime(v, anime); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Animes.Update(anime)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"anime": anime}, nil)
}

func (app *application) deleteAnimeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Animes.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "success"}, nil)
}
