package controller_html

import (
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/skrewby/yapper/models"
	"github.com/skrewby/yapper/types"
	views "github.com/skrewby/yapper/views/pages/threads"
)

func GetAllThreads(model models.Threads) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		threads, err := model.GetAllThreads()
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
		}

		templ.Handler(views.Threads(threads)).ServeHTTP(w, r)
	}
}

func GetAllThreadsStub(model models.Threads) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		threads, err := model.GetAllThreads()
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
		}

		templ.Handler(views.ThreadsStub(threads)).ServeHTTP(w, r)
	}
}

func NewThread() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(views.NewThread()).ServeHTTP(w, r)
	}
}

func NewThreadStub() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(views.NewThreadStub()).ServeHTTP(w, r)
	}
}

func CreateThread(model models.Threads) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.FormValue("title")
		ctx := r.Context()
		user, ok := ctx.Value("user").(*types.User)
		if !ok {
			http.Error(w, http.StatusText(422), 422)
			return
		}

		err := model.CreateThread(title, user.Id)
		if err != nil {
			slog.Error(err.Error())
			templ.Handler(views.NewThreadStubError(title)).ServeHTTP(w, r)
			http.Error(w, http.StatusText(500), 500)
			return
		}

		http.Redirect(w, r, "/threads/stub", 302)
	}
}
