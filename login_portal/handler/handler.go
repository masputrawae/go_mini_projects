package handler

import (
	"login_portal/model"
	"login_portal/repo"
	"login_portal/service"
	"login_portal/store"
	"login_portal/view"
	"net/http"
	"net/url"
	"path/filepath"
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

		if err := h.sessionStore.Set(w, id); err != nil {
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

		if err := h.sessionStore.Set(w, id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, path, http.StatusSeeOther)

	}

}

func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {

	session, ok := r.Context().Value("session").(model.Session)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, err := h.userService.GetUserByID(session.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := view.Dashboard(user, session.CSRFToken).Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *Handler) EditProfile(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:

		var p repo.UserUpdatePayload

		username := r.FormValue("username")
		password := r.FormValue("password")
		firstName := r.FormValue("first-name")
		lastName := r.FormValue("last-name")

		if username != "" {
			p.Username = new(username)
		}
		if password != "" {
			p.Password = new(password)
		}
		if firstName != "" {
			p.FirstName = new(firstName)
		}
		if lastName != "" {
			p.LastName = new(lastName)
		}

		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = h.userService.Edit(id, p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, filepath.Join("/user", strconv.Itoa(id)), http.StatusSeeOther)
	}
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	h.sessionStore.Delete(w, cookie.Value)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
