package controller_json

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/skrewby/yapper/auth"
	"github.com/skrewby/yapper/models"
	"github.com/skrewby/yapper/types"
)

func GetAllUsers(userModel models.Users) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := userModel.GetAllUsers()
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
		}

		res, _ := json.Marshal(user)
		w.Write(res)
	}
}

func CreateUser(userModel models.Users) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var user types.NewUser
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &user)

		hash, err := auth.HashPassword(user.Password)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}

		err = userModel.CreateUser(user.Email, user.Name, hash)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}

		w.Write([]byte("User created"))
	}
}

func GetUser(userModel models.Users) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user, ok := ctx.Value("user").(*types.User)
		if !ok {
			http.Error(w, http.StatusText(422), 422)
			return
		}

		res, _ := json.Marshal(user)
		w.Write(res)
	}
}

func UpdateUser(userModel models.Users) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user, ok := ctx.Value("user").(*types.User)
		if !ok {
			http.Error(w, http.StatusText(422), 422)
			return
		}

		var newUser types.User
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &newUser)
		if newUser.Name != "" {
			user.Name = newUser.Name
		}
		if newUser.Active != nil {
			user.Active = newUser.Active
		}
		if newUser.Password != "" {
			hash, err := auth.HashPassword(newUser.Password)
			if err != nil {
				http.Error(w, http.StatusText(500), 500)
				return
			}
			user.Password = hash
		}

		err := userModel.UpdateUser(user)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}

		user.Password = ""
		res, _ := json.Marshal(user)
		w.Write(res)
	}
}
