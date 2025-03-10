package controller

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/skrewby/yapper/auth"
	html "github.com/skrewby/yapper/controller/html"
	json "github.com/skrewby/yapper/controller/json"
	"github.com/skrewby/yapper/database"
	"github.com/skrewby/yapper/models"
	"github.com/skrewby/yapper/utils"
)

type Controller struct {
	jwt      *auth.JWT
	sessions *auth.Sessions
	models   Models
	env      utils.Environment
}

type Models struct {
	users   models.Users
	threads models.Threads
	editor  models.Editor
	files   models.Files
}

func NewController(env utils.Environment) Controller {
	jwt := auth.InitAuth(env.JWTSecret)
	db, err := database.ConnectDatabase(env)
	if err != nil {
		log.Fatal(err)
	}
	sessions := auth.InitSessionAuth()

	models := Models{
		users:   *models.NewUsersModel(db),
		threads: *models.NewThreadsModel(db),
		editor:  *models.NewEditorModel(db),
		files:   *models.NewFilesModel(db),
	}

	return Controller{
		jwt,
		sessions,
		models,
		env,
	}
}

func (c *Controller) StartServer() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	c.registerRoutes(r)

	http.ListenAndServe(":"+c.env.ServerPort, r)
}

func (c *Controller) registerRoutes(r chi.Router) {
	r.Mount("/", c.htmlRoutes())
	r.Mount("/api", c.jsonRoutes())
}

func (c *Controller) htmlRoutes() http.Handler {
	r := chi.NewRouter()

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(c.sessionAuth)

		setupFileServer(r)

		r.Get("/", html.Dashboard())

		r.Post("/markdown", html.ConvertMarkdownToHTML(c.models.editor))
		r.Post("/images", html.UploadImage(c.models.files))

		r.Route("/dashboard", func(r chi.Router) {
			r.Get("/", html.Dashboard())
			r.Get("/stub", html.DashboardStub())
		})

		r.Route("/threads", func(r chi.Router) {
			r.Get("/", html.GetAllThreads(c.models.threads))
			r.Get("/stub", html.GetAllThreadsStub(c.models.threads))

			r.Get("/new", html.NewThread())
			r.Get("/new/stub", html.NewThreadStub())
			r.Post("/new", html.CreateThread(c.models.threads))

		})

		r.Route("/settings", func(r chi.Router) {
			r.Get("/", html.Settings())
			r.Get("/stub", html.SettingsStub())
		})

		r.Route("/users", func(r chi.Router) {
			r.Get("/", html.GetAllUsers(c.models.users))
			r.Get("/stub", html.GetAllUsersStub(c.models.users))

			r.Get("/new", html.NewUser())
			r.Get("/new/stub", html.NewUserStub())
			r.Post("/new", html.CreateUser(c.models.users))

			r.Route("/{id}", func(r chi.Router) {
				r.Group(func(r chi.Router) {
					r.Use(c.usersCtx)

					r.Put("/", html.ChangeUserActiveStatus(c.models.users))
				})
			})

		})
	})

	// Public routes
	r.Group(func(r chi.Router) {
		r.Get("/login", html.LoginPage())
		r.Post("/login", html.Login(c.models.users, *c.sessions))
	})

	return r
}

func (c *Controller) jsonRoutes() http.Handler {
	r := chi.NewRouter()

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(c.jwt.Verify())
		r.Use(c.jwt.Authenticate())
		r.Use(c.setJsonHeader)

		r.Route("/users", func(r chi.Router) {
			r.Get("/", json.GetAllUsers(c.models.users))
			r.Post("/", json.CreateUser(c.models.users))

			r.Route("/{id}", func(r chi.Router) {
				r.Group(func(r chi.Router) {
					r.Use(c.usersCtx)

					r.Get("/", json.GetUser(c.models.users))
					r.Put("/", json.UpdateUser(c.models.users))
				})
			})
		})
	})

	// Public routes
	r.Group(func(r chi.Router) {
		r.Post("/login", json.Login(c.models.users, c.jwt))
		r.Post("/bootstrap", json.Bootstrap(c.models.users))
	})

	return r
}

func setupFileServer(r chi.Router) {
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, ".dev/files"))
	fileServer(r, "/files", filesDir)

}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
