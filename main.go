package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vishal2098govind/lenslocked/controllers"
	"github.com/vishal2098govind/lenslocked/models"
	userM "github.com/vishal2098govind/lenslocked/models/user"
	"github.com/vishal2098govind/lenslocked/templates"
	"github.com/vishal2098govind/lenslocked/views"
)

func main() {
	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}

	userS := userM.UserService{DB: db}
	usersC := controllers.Users{UserService: &userS}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get(
		"/",
		controllers.StaticHandler(views.Must(views.ParseFS(
			templates.FS,
			"home.gohtml",
			"layout-parts.gohtml",
		))),
	)

	r.Get(
		"/contact",
		controllers.StaticHandler(views.Must(views.ParseFS(
			templates.FS,
			"contact.gohtml",
			"layout-parts.gohtml",
		))),
	)

	r.Get(
		"/faq",
		controllers.FAQ(views.Must(views.ParseFS(
			templates.FS,
			"faq.gohtml",
			"layout-parts.gohtml",
		))),
	)

	usersC.Templates.New = views.Must(
		views.ParseFS(templates.FS, "signup.gohtml", "layout-parts.gohtml"),
	)

	r.Get("/signup", usersC.New)
	r.Post("/users", usersC.Create)

	r.Get("/products/{productId}", func(w http.ResponseWriter, r *http.Request) {
		pid := chi.URLParam(r, "productId")
		fmt.Fprint(w, "product id:", pid)
	})
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	fmt.Println("Starting the server on :3000")
	http.ListenAndServe(":3000", r)
}
