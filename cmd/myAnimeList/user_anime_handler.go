package main

import (
	"errors"
	"log"
	"net/http"

	"final-project/pkg/model"
	"final-project/pkg/validator"

	"github.com/lib/pq"
)

func (app *application) getAverageRatingHandler(w http.ResponseWriter, r *http.Request) {
	animeID, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	avgRating, userCount, title, err := app.models.User_and_Anime.AverageRating(animeID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	envelope := envelope{
		"title":      title,
		"animeID":    animeID,
		"avg_rating": avgRating,
		"user_count": userCount,
	}

	err = app.writeJSON(w, http.StatusOK, envelope, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getUserAnimesByUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		model.Filters
	}
	v := validator.New()
	qs := r.URL.Query()

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readStrings(qs, "sort", "id")
	input.Filters.SortSafeList = []string{
		"id", "-id",
		"rating", "-rating",
		"created_at", "-created_at",
		"updated_at", "-updated_at",
	}

	if model.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	userAnimes, metadata, err := app.models.User_and_Anime.GetAllByUser(userID, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"userAnimes": userAnimes, "metadata": metadata}, nil)
}

func (app *application) createUserAnimeHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID  int     `json:"userID"`
		AnimeID int     `json:"animeID"`
		Rating  float64 `json:"rating"`
		Review  string  `json:"review"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		log.Println(err)
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}

	userAnime := &model.User_and_Anime{
		UserID:  input.UserID,
		AnimeID: input.AnimeID,
		Rating:  input.Rating,
		Review:  input.Review,
	}

	v := validator.New()
	if model.ValidateUA(v, userAnime); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.User_and_Anime.Insert(userAnime)
	if err != nil {
		log.Println("Database error:", err) // Log the error
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505": // Unique violation
				app.errorResponse(w, r, http.StatusConflict, "UserAnime already exists")
				return
			case "23503": // Foreign key violation
				if pqErr.Constraint == "user_and_anime_userid_fkey" {
					app.errorResponse(w, r, http.StatusBadRequest, "Invalid userID")
					return
				}
				if pqErr.Constraint == "user_and_anime_animeid_fkey" {
					app.errorResponse(w, r, http.StatusBadRequest, "Invalid animeID")
					return
				}
			case "23502": // Not null violation
				app.errorResponse(w, r, http.StatusBadRequest, "A required field is missing")
				return
			default: // General database error
				app.serverErrorResponse(w, r, err)
				return
			}
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	app.writeJSON(w, http.StatusCreated, envelope{"userAnime": userAnime}, nil)
}

func (app *application) getUserAnimeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	userAnime, err := app.models.User_and_Anime.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"userAnime": userAnime}, nil)
}

func (app *application) updateUserAnimeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	userAnime, err := app.models.User_and_Anime.Get(id)
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
		Review *string  `json:"review"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Rating != nil {
		userAnime.Rating = *input.Rating
	}
	if input.Review != nil {
		userAnime.Review = *input.Review
	}

	v := validator.New()
	if model.ValidateUA(v, userAnime); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.User_and_Anime.Update(*userAnime)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"userAnime": userAnime}, nil)
}

func (app *application) deleteUserAnimeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.User_and_Anime.Delete(id)
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
