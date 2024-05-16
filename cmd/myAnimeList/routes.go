package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// routes is our main application's router.
func (app *application) routes() http.Handler {
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(app.notFoundResponse)

	// Convert app.methodNotAllowedResponse helper to a http.Handler and set it as the custom
	// error handler for 405 Method Not Allowed responses
	r.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowedResponse)

	// healthcheck
	r.HandleFunc("/api/v1/healthcheck", app.healthcheckHandler).Methods("GET")

	anime1 := r.PathPrefix("/api/v1").Subrouter()

	// Menu Singleton
	// localhost:8081/api/v1/animes
	anime1.HandleFunc("/animes", app.getAnimesList).Methods("GET")
	anime1.HandleFunc("/animes", app.requirePermissions("animes:write", app.createAnimeHandler)).Methods("POST")
	anime1.HandleFunc("/animes/{id:[0-9]+}", app.getAnimeHandler).Methods("GET")
	anime1.HandleFunc("/animes/{id:[0-9]+}", app.requirePermissions("animes:write", app.updateAnimeHandler)).Methods("PUT")
	anime1.HandleFunc("/animes/{id:[0-9]+}", app.requirePermissions("animes:write", app.deleteAnimeHandler)).Methods("DELETE")

	users1 := r.PathPrefix("/api/v1").Subrouter()
	// User handlers with Authentication
	users1.HandleFunc("/users", app.registerUserHandler).Methods("POST") // TODO: the value gets into DB, but server doesn't respond + log.error
	users1.HandleFunc("/users/activated", app.activateUserHandler).Methods("PUT")
	users1.HandleFunc("/users/login", app.createAuthenticationTokenHandler).Methods("POST")

	users_animes1 := r.PathPrefix("/api/v1").Subrouter()
	// User and Anime handlers
	users_animes1.HandleFunc("/user_animes", app.createUserAnimeHandler).Methods("POST")
	users_animes1.HandleFunc("/user_animes/user/{id:[0-9]+}", app.getUserAnimesByUserHandler).Methods("GET")
	users_animes1.HandleFunc("/user_animes/{id:[0-9]+}", app.getUserAnimeHandler).Methods("GET")
	users_animes1.HandleFunc("/user_animes/{id:[0-9]+}", app.updateUserAnimeHandler).Methods("PUT")
	users_animes1.HandleFunc("/user_animes/{id:[0-9]+}", app.deleteUserAnimeHandler).Methods("DELETE")

	// AVG rating by users
	users_animes1.HandleFunc("/user_animes/animes/{id:[0-9]+}", app.getAverageRatingHandler).Methods("GET")

	return app.authenticate(r)
}
