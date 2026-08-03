package main

import (
	"log"
	"login_portal/handler"
	"login_portal/middleware"
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
	mw := middleware.New(sessionStore)
	handler := handler.New(userService, sessionStore)

	mux.HandleFunc("/", handler.Root)
	mux.HandleFunc("/login", handler.Login)
	mux.HandleFunc("/register", handler.Register)
	mux.HandleFunc("/user/{id}",
		mw.Auth(
			mw.CSRF(
				mw.UserPath(
					handler.Dashboard,
				),
			),
		),
	)

	mux.HandleFunc("/user/{id}/edit",
		mw.Auth(
			mw.CSRF(
				mw.UserPath(
					handler.EditProfile,
				),
			),
		),
	)

	mux.HandleFunc("/logout", handler.Logout)
	log.Fatal(http.ListenAndServe(":8080", mw.Logger(mux)))
}
