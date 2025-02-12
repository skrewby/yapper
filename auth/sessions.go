package auth

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type sessionToken struct {
	Token  string
	Expiry time.Time
}

type session struct {
	email  string
	expiry time.Time
}

type Sessions struct {
	sessions map[string]session
}

func InitSessionAuth() *Sessions {
	return &Sessions{
		sessions: map[string]session{},
	}
}

func (s *Sessions) Create(w http.ResponseWriter, email string, remember bool) {
	token := uuid.NewString()
	expiresAt := time.Now().Add(60 * 60 * 24 * 30 * time.Second)

	s.sessions[token] = session{
		email:  email,
		expiry: expiresAt,
	}

	t := sessionToken{
		token,
		expiresAt,
	}
	s.setCookie(w, t, remember)
}

func (s *Sessions) Valid(token string) bool {
	sess, exists := s.sessions[token]
	if !exists {
		slog.Info("Attempted to log in with token that does not exist")
		return false
	}
	if sess.isExpired() {
		slog.Info("Token expired")
		delete(s.sessions, token)
		return false
	}

	return true
}

func (s *Sessions) Delete(w http.ResponseWriter, token string) {
	_, exists := s.sessions[token]
	if !exists {
		slog.Info("Attempted to delete session that does not exist")
		return
	}

	delete(s.sessions, token)
	s.clearCookie(w)
}

func (s *Sessions) setCookie(w http.ResponseWriter, t sessionToken, remember bool) {
	if remember {
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    t.Token,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
			HttpOnly: true,
			Expires:  t.Expiry,
		})
	} else {
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    t.Token,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
			HttpOnly: true,
		})

	}
}

func (s *Sessions) clearCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Expires:  time.Now(),
	})
}

func (s *session) isExpired() bool {
	return s.expiry.Before(time.Now())
}
