package middleware

import (
	"context"
	"log"
	"login_portal/model"
	"login_portal/store"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Middleware struct {
	sessionStore *store.Session
}

func New(sessionStore *store.Session) *Middleware {
	return &Middleware{sessionStore: sessionStore}
}

func (m *Middleware) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		log.Println(r.Method, r.URL.Path, time.Since(start))
	})
}

func (m *Middleware) Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("session")
		if err != nil {

			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		sessionID := cookie.Value
		session, err := m.sessionStore.Get(sessionID)
		if err != nil {

			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if session.ExpiresAt.Compare(time.Now()) <= 0 {
			m.sessionStore.Delete(w, sessionID)

			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), "session", session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) CSRF(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodTrace:
			next.ServeHTTP(w, r)
		default:
			session, ok := r.Context().Value("session").(model.Session)
			if ok {

				token := r.FormValue("csrf-token")
				if token != session.CSRFToken {
					http.Error(w, "token csrf not matches", http.StatusForbidden)
					return
				}

				next.ServeHTTP(w, r)

			}
		}

	})
}

func (m *Middleware) UserPath(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("session")
		if err != nil {

			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return

		}

		sessionID := cookie.Value

		session, err := m.sessionStore.Get(sessionID)
		if err != nil {

			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return

		}

		id, err := strconv.Atoi(r.PathValue("id"))
		if session.UserID != id {

			path, err := url.JoinPath("/user", strconv.Itoa(session.UserID))
			if err != nil {

				http.Error(w, err.Error(), http.StatusInternalServerError)
				return

			}

			http.Redirect(w, r, path, http.StatusSeeOther)
			return

		}

		next.ServeHTTP(w, r)
	})
}
