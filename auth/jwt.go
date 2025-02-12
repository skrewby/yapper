package auth

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/skrewby/yapper/types"
)

type JWT struct {
	auth *jwtauth.JWTAuth
}

func InitAuth(secret string) *JWT {
	auth := jwtauth.New("HS256", []byte(secret), nil)

	return &JWT{
		auth,
	}
}

func (jwt JWT) CreateToken(context types.JWTContext) string {
	_, tokenString, _ := jwt.auth.Encode(map[string]interface{}{
		"context": map[string]interface{}{
			"user": map[string]interface{}{
				"email":        context.User.Email,
				"display_name": context.User.Name,
			},
		},
	})

	return tokenString
}

func (jwt JWT) Verify() func(http.Handler) http.Handler {
	return jwtauth.Verifier(jwt.auth)
}

func (jwt JWT) Authenticate() func(http.Handler) http.Handler {
	return jwtauth.Authenticator(jwt.auth)
}

func (jwt JWT) Token(r *http.Request) types.JWT {
	_, token, _ := jwtauth.FromContext(r.Context())

	var ctx types.JWT
	d, _ := json.Marshal(token)
	_ = json.Unmarshal(d, &ctx)

	return ctx
}
