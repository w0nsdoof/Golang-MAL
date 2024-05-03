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
	// localhost:8081/api/v1/menus
	anime1.HandleFunc("/animes", app.getAnimesList).Methods("GET") // TODO:
	anime1.HandleFunc("/animes", app.createAnimeHandler).Methods("POST")
	anime1.HandleFunc("/animes/{id:[0-9]+}", app.getAnimeHandler).Methods("GET")
	anime1.HandleFunc("/animes/{id:[0-9]+}", app.updateAnimeHandler).Methods("PUT")
	anime1.HandleFunc("/animes/{id:[0-9]+}", app.requirePermissions("animes:write", app.deleteAnimeHandler)).Methods("DELETE")

	users1 := r.PathPrefix("/api/v1").Subrouter()
	// User handlers with Authentication
	users1.HandleFunc("/users", app.registerUserHandler).Methods("POST")
	users1.HandleFunc("/users/activated", app.activateUserHandler).Methods("PUT")
	users1.HandleFunc("/users/login", app.createAuthenticationTokenHandler).Methods("POST")

	return app.authenticate(r)
}
