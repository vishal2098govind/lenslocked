package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/vishal2098govind/lenslocked/controllers"
	"github.com/vishal2098govind/lenslocked/middlewares"
	"github.com/vishal2098govind/lenslocked/migrations"
	postgresDB "github.com/vishal2098govind/lenslocked/models/db"
	sessionM "github.com/vishal2098govind/lenslocked/models/session"
	userM "github.com/vishal2098govind/lenslocked/models/user"
	"github.com/vishal2098govind/lenslocked/templates"
	"github.com/vishal2098govind/lenslocked/views"
)

func main() {
	// setup database
	cfg := postgresDB.DefaultPostgresConfig()
	db, err := postgresDB.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = postgresDB.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// setup services
	userS := userM.UserService{DB: db}
	sessionS := sessionM.SessionService{DB: db}

	// setup controllers
	usersC := controllers.Users{
		UserService:    &userS,
		SessionService: &sessionS,
	}
	usersC.Templates.New = views.Must(
		views.ParseFS(templates.FS, "signup.gohtml", "layout-parts.gohtml"),
	)
	usersC.Templates.SignIn = views.Must(
		views.ParseFS(templates.FS, "signin.gohtml", "layout-parts.gohtml"),
	)
	usersC.Templates.CurrentUser = views.Must(
		views.ParseFS(templates.FS, "current_user.gohtml", "layout-parts.gohtml"),
	)

	// setup middlewares
	usersMW := middlewares.UsersMiddleware{SessionService: &sessionS}
	csrfKey := "gFvi45R4fy5xNBlnEeZtQbfAVCYEIAUX"
	csrfMw := csrf.Protect(
		[]byte(csrfKey),
		// TODO: Fix this before deploying. requires HTTPS
		csrf.Secure(false),
	)

	// setup router and routes
	r := chi.NewRouter()
	r.Use(csrfMw)
	r.Use(middlewares.IPLoggerMiddleware)
	r.Use(usersMW.SetUser)
	r.Get("/", controllers.StaticHandler(views.Must(views.ParseFS(
		templates.FS,
		"home.gohtml",
		"layout-parts.gohtml",
	))))
	r.Get("/contact", controllers.StaticHandler(views.Must(views.ParseFS(
		templates.FS,
		"contact.gohtml",
		"layout-parts.gohtml",
	))))
	r.Get("/faq", controllers.FAQ(views.Must(views.ParseFS(
		templates.FS,
		"faq.gohtml",
		"layout-parts.gohtml",
	))))
	r.Get("/signup", usersC.New)
	r.Get("/signin", usersC.SignIn)
	r.Post("/users", usersC.Create)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Post("/signout", usersC.ProcessSignOut)
	r.Route("/users/me", func(r chi.Router) {
		r.Use(usersMW.RequireUser)
		r.Get("/", usersC.CurrentUser)
	})

	r.Get("/products/{productId}", func(w http.ResponseWriter, r *http.Request) {
		pid := chi.URLParam(r, "productId")
		fmt.Fprint(w, "product id:", pid)
	})
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	// serve HTTP
	fmt.Println("Starting the server on :3000")
	err = http.ListenAndServe(":3000", r)
	if err != nil {
		panic(err)
	}
}
