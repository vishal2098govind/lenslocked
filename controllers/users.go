package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/vishal2098govind/lenslocked/context"
	"github.com/vishal2098govind/lenslocked/cookies"
	sessionM "github.com/vishal2098govind/lenslocked/models/session"
	userM "github.com/vishal2098govind/lenslocked/models/user"
)

type Users struct {
	UserService *userM.UserService

	SessionService *sessionM.SessionService

	Templates struct {
		New         Template // signup
		SignIn      Template // signin
		CurrentUser Template // current user
	}
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	if context.User(r.Context()) != nil {
		http.Redirect(w, r, "/users/me", http.StatusFound)
		return
	}
	var data struct{ Email string }
	data.Email = r.FormValue("email")

	u.Templates.New.Execute(w, r, data)
}

func (u Users) SignIn(w http.ResponseWriter, r *http.Request) {
	if context.User(r.Context()) != nil {
		http.Redirect(w, r, "/users/me", http.StatusFound)
		return
	}

	var data struct{ Email string }
	data.Email = r.FormValue("email")

	u.Templates.SignIn.Execute(w, r, data)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	res, err := u.UserService.Create(userM.CreateUserRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	sessRes, err := u.SessionService.Create(sessionM.CreateSessionRequest{
		UserID: res.User.ID,
	})
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	cookies.SetCookie(w, cookies.CookieSession, sessRes.Session.Token)

	http.Redirect(w, r, "/users/me", http.StatusFound)
}

func (u Users) ProcessSignIn(w http.ResponseWriter, r *http.Request) {
	email := strings.ToLower(r.FormValue("email"))
	password := r.FormValue("password")

	res, err := u.UserService.Authenticate(userM.AuthenticateRequest{
		Email:    email,
		Password: password,
	})
	if err == userM.ErrUserNotFound {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}
	if err == userM.ErrInvalidCredentials {
		http.Error(w, "Invalid email or password", http.StatusBadRequest)
		return
	}
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	sessRes, err := u.SessionService.Create(sessionM.CreateSessionRequest{
		UserID: res.User.ID,
	})
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	cookies.SetCookie(w, cookies.CookieSession, sessRes.Session.Token)

	http.Redirect(w, r, "/users/me", http.StatusFound)

}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := context.Logger(ctx)
	logger.Printf("current user")

	user := context.User(ctx)

	u.Templates.CurrentUser.Execute(w, r, user)
}

func (u Users) ProcessSignOut(w http.ResponseWriter, r *http.Request) {
	sessionC, err := cookies.ReadCookie(r, cookies.CookieSession)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	_, err = u.SessionService.Delete(&sessionM.DeleteSessionRequest{
		Token: sessionC,
	})
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	cookies.DeleteCookie(w, cookies.CookieSession)
	http.Redirect(w, r, "/signin", http.StatusFound)
}
