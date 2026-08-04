package main

import (
	"database/sql"
	"log"
	"net/http"
	"simple_todo/handler"
	"simple_todo/middleware"
	"simple_todo/repo"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	mux := http.NewServeMux()
	repo := repo.New(ConnectDB())
	mw := middleware.New()
	handler := handler.New(repo)

	mux.HandleFunc("GET /", handler.GetAll)
	mux.HandleFunc("GET /{id}", handler.GetByID)
	mux.HandleFunc("POST /add", handler.Add)
	mux.HandleFunc("DELETE /{id}/delete", handler.Delete)
	mux.HandleFunc("PATCH /{id}/edit", handler.Update)

	m := mw.ContentTypeJson(mux)
	m = mw.Logger(m)

	http.ListenAndServe(":8080", m)
}

func ConnectDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./app.db")
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}
