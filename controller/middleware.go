package controller

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (c *Controller) setJsonHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func (c *Controller) sessionAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				slog.Error("Session Auth: No cookie")

				http.Redirect(w, r, "/login", 302)
				return
			}

			slog.Error(err.Error())
			http.Redirect(w, r, "/login", 302)
			return
		}
		token := cookie.Value

		if !c.sessions.Valid(token) {
			slog.Error("Session Auth: Session not valid")

			http.Redirect(w, r, "/login", 302)
			return
		}

		userId, err := c.sessions.GetUserId(token)
		if err != nil {
			http.Redirect(w, r, "/login", 302)
			return
		}

		user, err := c.models.users.GetUser(userId)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		ctx := context.WithValue(r.Context(), "user", user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (c *Controller) usersCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		user, err := c.models.users.GetUser(id)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
