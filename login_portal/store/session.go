package store

import (
	"errors"
	"login_portal/model"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	ErrSessionIDNotExist = errors.New("session id not exist")
)

type Session struct {
	sessions map[string]model.Session
	mu       sync.RWMutex
}

func New() *Session {
	return &Session{sessions: make(map[string]model.Session)}
}

func (s *Session) Get(sessionID string) (model.Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// temukan session id
	session, exist := s.sessions[sessionID]
	if !exist {
		return model.Session{}, ErrSessionIDNotExist
	}

	return session, nil
}

func (s *Session) Set(w http.ResponseWriter, userID int) error {
	// buat token untuk session
	uuidSession, err := uuid.NewV7()
	if err != nil {
		return err
	}

	// buat token untuk csrf
	uuidCSRF, err := uuid.NewV7()
	if err != nil {
		return err
	}

	// tanggal kadaluarsa
	expires := time.Now().Add(24 * time.Hour)

	s.mu.Lock()
	s.sessions[uuidSession.String()] = model.Session{
		UserID:    userID,
		CSRFToken: uuidCSRF.String(),
		ExpiresAt: expires,
	}
	s.mu.Unlock()

	// simpan session id di cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    uuidSession.String(),
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		MaxAge:   int(24 * time.Hour),
		SameSite: http.SameSiteLaxMode,
	})

	return nil
}

func (s *Session) Delete(w http.ResponseWriter, sessionID string) {
	s.mu.Lock()

	// hapus session di memori
	delete(s.sessions, sessionID)

	s.mu.Unlock()

	// hapus session id di cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		MaxAge:   -1,
		SameSite: http.SameSiteLaxMode,
	})
}
