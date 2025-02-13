package controller_html

import (
	"log/slog"
	"net/http"

	"github.com/a-h/templ"

	"github.com/skrewby/yapper/auth"
	"github.com/skrewby/yapper/models"
	"github.com/skrewby/yapper/views/pages"
)

func LoginPage() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(views.Login(views.LoginErrorNoError)).ServeHTTP(w, r)
	}
}

func Login(userModel models.Users, sessions auth.Sessions) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.FormValue("email")
		password := r.FormValue("password")
		remember := r.FormValue("remember")

		user, err := userModel.GetUserByEmail(email)
		if err != nil {
			slog.Error("Login - could not find user")
			templ.Handler(views.Login(views.LoginErrorInvalidCredentials)).ServeHTTP(w, r)
			return
		}

		hash, err := userModel.GetUserHashedPassword(user.Id)
		if err != nil {
			slog.Error("Login - could not find user hashed password")
			templ.Handler(views.Login(views.LoginErrorInvalidCredentials)).ServeHTTP(w, r)
			return
		}

		valid := auth.ValidPassword(password, *hash)
		if !valid {
			slog.Error("Login - not a valid password")
			templ.Handler(views.Login(views.LoginErrorInvalidCredentials)).ServeHTTP(w, r)
			return
		}

		sessions.Create(w, user.Id, user.Email, remember == "on")
		http.Redirect(w, r, "/", 302)
	}
}

func Logout(sessions auth.Sessions) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Redirect(w, r, "/login", 302)
				return
			}
			slog.Error(err.Error())
			http.Redirect(w, r, "/login", 302)
			return
		}
		token := cookie.Value

		sessions.Delete(w, token)

		templ.Handler(views.Login(views.LoginErrorNoError)).ServeHTTP(w, r)
	}
}
