package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/joho/godotenv"
	"github.com/vishal2098govind/lenslocked/controllers"
	"github.com/vishal2098govind/lenslocked/middlewares"
	"github.com/vishal2098govind/lenslocked/migrations"
	postgresDB "github.com/vishal2098govind/lenslocked/models/db"
	emailM "github.com/vishal2098govind/lenslocked/models/email"
	passwordResetM "github.com/vishal2098govind/lenslocked/models/password_reset"
	sessionM "github.com/vishal2098govind/lenslocked/models/session"
	userM "github.com/vishal2098govind/lenslocked/models/user"
	"github.com/vishal2098govind/lenslocked/templates"
	"github.com/vishal2098govind/lenslocked/views"
)

type config struct {
	PSQL postgresDB.PostgresConfig
	SMTP emailM.SMTPConfig
	CSRF struct {
		Key    string
		Secure bool
	}
	Server struct {
		Address string
	}
}

func LoadEnvConfig() (config, error) {
	var cfg config

	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}

	cfg.PSQL = postgresDB.DefaultPostgresConfig()

	cfg.CSRF.Key = "gFvi45R4fy5xNBlnEeZtQbfAVCYEIAUX"
	cfg.CSRF.Secure = false

	cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	cfg.SMTP.Port, err = strconv.Atoi(portStr)
	if err != nil {
		return cfg, err
	}
	cfg.SMTP.Username = os.Getenv("SMTP_USERNAME")
	cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")

	cfg.Server.Address = ":3000"

	return cfg, nil
}

func main() {
	cfg, err := LoadEnvConfig()
	if err != nil {
		panic(err)
	}

	// setup database
	db, err := postgresDB.Open(cfg.PSQL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = postgresDB.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// setup services
	userS := &userM.UserService{DB: db}
	sessionS := &sessionM.SessionService{DB: db}
	emailS := emailM.NewEmailService(cfg.SMTP)
	passwordResetS := &passwordResetM.PasswordResetService{DB: db}

	// setup controllers
	usersC := controllers.Users{
		UserService:          userS,
		SessionService:       sessionS,
		PasswordResetService: passwordResetS,
		EmailService:         emailS,
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
	usersC.Templates.ForgotPassword = views.Must(
		views.ParseFS(templates.FS, "forgot_password.html", "layout-parts.gohtml"),
	)
	usersC.Templates.ResetPassword = views.Must(
		views.ParseFS(templates.FS, "reset_password.html", "layout-parts.gohtml"),
	)

	emailS.Templates.ForgotPasswordTpl = template.Must(
		template.ParseFS(templates.FS, "forgot_password_email.gohtml"),
	)

	// setup middlewares
	loggerMW := middlewares.LoggerMiddleware{}
	usersMW := middlewares.UsersMiddleware{SessionService: sessionS}
	csrfMw := csrf.Protect(
		[]byte(cfg.CSRF.Key),
		csrf.Secure(cfg.CSRF.Secure),
	)

	// setup router and routes
	r := chi.NewRouter()
	r.Use(csrfMw)
	r.Use(loggerMW.IPLoggerMiddleware)
	r.Use(usersMW.SetUser)
	r.Use(loggerMW.Logger)
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
	r.Get("/forgot-password", usersC.ForgotPassword)
	r.Post("/forgot-password", usersC.ProcessForgotPassword)
	r.Get("/reset-password", usersC.ResetPassword)
	r.Post("/reset-password", usersC.ProcessResetPassword)

	r.Get("/products/{productId}", func(w http.ResponseWriter, r *http.Request) {
		pid := chi.URLParam(r, "productId")
		fmt.Fprint(w, "product id:", pid)
	})
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	// serve HTTP
	fmt.Printf("Starting the server on %s...\n", cfg.Server.Address)
	err = http.ListenAndServe(cfg.Server.Address, r)
	if err != nil {
		panic(err)
	}
}
