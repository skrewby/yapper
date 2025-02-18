package controller_html

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/skrewby/yapper/models"
	views "github.com/skrewby/yapper/views/pages"
)

func ConvertMarkdownToHTML(model models.Editor) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		md := r.FormValue("message")
		html, err := model.ConvertMarkdownToHTML([]byte(md))
		if err != nil {
			templ.Handler(views.PreviewError()).ServeHTTP(w, r)
			return
		}

		templ.Handler(views.Preview(html)).ServeHTTP(w, r)
	}
}

