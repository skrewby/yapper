package controller_html

import (
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/skrewby/yapper/auth"
	"github.com/skrewby/yapper/models"
	"github.com/skrewby/yapper/types"
	"github.com/skrewby/yapper/utils"
	views "github.com/skrewby/yapper/views/pages/users"
)

func GetAllUsers(userModel models.Users) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := userModel.GetAllUsers()
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
		}

		templ.Handler(views.Users(users)).ServeHTTP(w, r)
	}
}

func GetAllUsersStub(userModel models.Users) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := userModel.GetAllUsers()
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
		}

		templ.Handler(views.UsersStub(users)).ServeHTTP(w, r)
	}
}

func ChangeUserActiveStatus(userModel models.Users) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user, ok := ctx.Value("user").(*types.User)
		if !ok {
			http.Error(w, http.StatusText(422), 422)
			return
		}

		active := r.FormValue("active")
		r.FormValue("reset_password")

		if active == "true" {
			user.Active = utils.Pointer(true)
		} else if active == "false" {
			user.Active = utils.Pointer(false)
		}

		err := userModel.UpdateUser(user)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}

		users, err := userModel.GetAllUsers()
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
		}

		templ.Handler(views.UsersStub(users)).ServeHTTP(w, r)
	}
}

func NewUser() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(views.NewUser()).ServeHTTP(w, r)
	}
}

func NewUserStub() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(views.NewUserStub()).ServeHTTP(w, r)
	}
}

func CreateUser(userModel models.Users) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.FormValue("email")
		displayName := r.FormValue("displayName")
		password := r.FormValue("password")

		hash, err := auth.HashPassword(password)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}

		err = userModel.CreateUser(email, displayName, hash)
		if err != nil {
			slog.Error(err.Error())
			templ.Handler(views.NewUserStubError(email, displayName, password, err)).ServeHTTP(w, r)
			http.Error(w, http.StatusText(500), 500)
			return
		}

		http.Redirect(w, r, "/users/stub", 302)
	}
}
