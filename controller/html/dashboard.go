package controller_html

import (
	"net/http"

	"github.com/a-h/templ"
	views "github.com/skrewby/yapper/views/pages"
)

func Dashboard() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(views.Dashboard()).ServeHTTP(w, r)
	}
}

func DashboardStub() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(views.DashboardStub()).ServeHTTP(w, r)
	}
}
