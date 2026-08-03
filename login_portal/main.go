package main

import (
	"log"
	"login_portal/handler"
	"login_portal/model"
	"login_portal/repo"
	"login_portal/service"
	"login_portal/store"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	userRepo := repo.NewUser([]model.User{})
	userService := service.NewUser(userRepo)
	sessionStore := store.New()
	handler := handler.New(userService, sessionStore)

	mux.HandleFunc("/", handler.Root)
	mux.HandleFunc("/login", handler.Login)
	mux.HandleFunc("/register", handler.Register)
	mux.HandleFunc("/user/{id}", handler.Dashboard)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
