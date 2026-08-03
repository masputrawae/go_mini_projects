package handler

import (
	"login_portal/service"
	"login_portal/store"
	"login_portal/view"
	"net/http"
	"net/url"
	"strconv"
)

type Handler struct {
	userService  *service.User
	sessionStore *store.Session
}

func New(userService *service.User, sessionStore *store.Session) *Handler {
	return &Handler{userService: userService, sessionStore: sessionStore}
}

func (h *Handler) Root(w http.ResponseWriter, r *http.Request) {
	if err := view.Index(view.Guest()).Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if err := view.Login(nil).Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case http.MethodPost:

		username := r.FormValue("username")
		password := r.FormValue("password")

		id, err := h.userService.Login(username, password)
		if err != nil {

			if err := view.Login(err).Render(r.Context(), w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			return
		}

		path, err := url.JoinPath("/user", strconv.Itoa(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, path, http.StatusSeeOther)
	}

}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:

		if err := view.Register(nil).Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case http.MethodPost:

		username := r.FormValue("username")
		password := r.FormValue("password")

		id, err := h.userService.Register(username, password)
		if err != nil {

			if err := view.Register(err).Render(r.Context(), w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			return
		}

		path, err := url.JoinPath("/user", strconv.Itoa(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, path, http.StatusSeeOther)

	}

}

func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := view.Dashboard(user).Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {}
