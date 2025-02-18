package controller_html

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/skrewby/yapper/models"
	views "github.com/skrewby/yapper/views/components"
)

func UploadImage(model models.Files) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Max size of files is 10MB
		r.ParseMultipartForm(10 << 20)

		file, _, err := r.FormFile("file")
		if err != nil {
			templ.Handler(views.ImageUploadError()).ServeHTTP(w, r)
			return
		}
		defer file.Close()

		path, err := model.UploadImage(file)
		if err != nil {
			templ.Handler(views.ImageUploadError()).ServeHTTP(w, r)
			return
		}

		templ.Handler(views.Image(path)).ServeHTTP(w, r)
	}
}
