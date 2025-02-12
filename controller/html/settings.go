package controller_html

import (
	"net/http"

	"github.com/a-h/templ"
	views "github.com/skrewby/yapper/views/pages"
)

func Settings() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(views.Settings()).ServeHTTP(w, r)
	}
}

func SettingsStub() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(views.SettingsStub()).ServeHTTP(w, r)
	}
}
