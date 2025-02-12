package controller_json

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/skrewby/yapper/auth"
	"github.com/skrewby/yapper/models"
	"github.com/skrewby/yapper/types"
)

func Login(userModel models.Users, jwt *auth.JWT) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var login types.Login
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &login)

		user, err := userModel.GetUserByEmail(login.Email)
		if err != nil {
			http.Error(w, "Username or password incorrect", 401)
			return
		}
		hash, err := userModel.GetUserHashedPassword(user.Id)
		if err != nil {
			http.Error(w, "Username or password incorrect", 401)
			return
		}
		valid := auth.ValidPassword(login.Password, *hash)
		if !valid {
			http.Error(w, "Username or password incorrect", 401)
			return
		}

		usrCtx := types.JWTUser{
			Email: user.Email,
			Name:  user.Name,
		}
		ctx := types.JWTContext{
			User: usrCtx,
		}

		token := jwt.CreateToken(ctx)
		res, _ := json.Marshal(&types.Auth{
			Token: token,
		})
		w.Write(res)

	}
}

func Bootstrap(userModel models.Users) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := userModel.GetAllUsers()
		if users != nil {
			http.Error(w, "Already bootstrapped", 403)
			return
		}
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
		}

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

		w.Write([]byte("Success"))
	}
}
