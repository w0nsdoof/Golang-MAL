package main

import (
	"database/sql"
	"final-project/pkg/models"
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

type config struct {
	port string
	env  string
	db   struct {
		dsn string
	}
}

type application struct {
	config    config
	userModel *models.UserModel
}

func main() {
	var cfg config
	flag.StringVar(&cfg.port, "port", ":8081", "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgresql://postgres:1473@localhost/golang?sslmode=disable", "PostgreSQL DSN")
	flag.Parse()

	db, err := openDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userModel := &models.UserModel{DB: db}

	app := &application{
		config:    cfg,
		userModel: userModel,
	}

	app.run()
}

func (app *application) run() {
	r := mux.NewRouter()

	v1 := r.PathPrefix("/api/v1").Subrouter()

	userHandler := &UserHandler{Model: app.userModel}

	v1.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	v1.HandleFunc("/users/{id:[0-9]+}", userHandler.GetUser).Methods("GET")
	v1.HandleFunc("/users/{id:[0-9]+}", userHandler.UpdateUser).Methods("PUT")
	v1.HandleFunc("/users/{id:[0-9]+}", userHandler.DeleteUser).Methods("DELETE")

	v1.HandleFunc("/users", userHandler.GetAllUsers).Methods("GET")

	log.Printf("Starting server on %s\n", app.config.port)
	err := http.ListenAndServe(app.config.port, r)
	log.Fatal(err)
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
